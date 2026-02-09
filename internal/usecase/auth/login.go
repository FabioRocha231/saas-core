package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type LoginUsecase struct {
	userRepo     repository.UserRepository
	storeRepo    repository.StoreRepository
	sessionRepo  repository.SessionRepository
	jwtService   ports.JwtInterface
	passwordHash ports.PasswordHashInterface
	context      context.Context
}

type LoginInput struct {
	Email    string
	Password string
}

type UserLoginOutput struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  string          `json:"name"`
	Role  entity.UserKind `json:"role"`
}

type LoginOutput struct {
	Token       string          `json:"token"`
	User        UserLoginOutput `json:"user"`
	StoresCount int             `json:"stores_count"`
	NextStep    entity.NextStep `json:"next_step"`
}

func NewLoginUsecase(
	context context.Context,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	storeRepo repository.StoreRepository,
	jwtService ports.JwtInterface,
	passwordHash ports.PasswordHashInterface,
) *LoginUsecase {
	return &LoginUsecase{
		context:      context,
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		storeRepo:    storeRepo,
		jwtService:   jwtService,
		passwordHash: passwordHash,
	}
}

func (l *LoginUsecase) Execute(input LoginInput) (*LoginOutput, error) {
	user, err := l.userRepo.GetByMail(l.context, input.Email)

	if err != nil {
		return nil, err
	}

	isEqual := l.passwordHash.Verify(user.Password, input.Password)
	if !isEqual {
		return nil, errx.New(errx.CodeNotFound, "user not found")
	}

	token, err := l.jwtService.Sign(user.ID, user.Role.String())
	if err != nil {
		return nil, err
	}

	err = l.sessionRepo.Create(l.context, &entity.Session{
		ID:        l.jwtService.GetJTI(token),
		UserID:    user.ID,
		ExpiresAt: l.jwtService.GetExpiresAt(token),
		Role:      user.Role.String(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	storesQuantity, err := l.storeRepo.CountByOwnerID(l.context, user.ID)
	if err != nil {
		return nil, err
	}

	userRole := mapRoleToKind(user.Role)

	return &LoginOutput{
		Token: token,
		User: UserLoginOutput{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  userRole,
		},
		StoresCount: storesQuantity,
		NextStep:    decideNextStep(userRole, storesQuantity),
	}, nil
}

func mapRoleToKind(role entity.UserRole) entity.UserKind {
	switch role {
	case entity.UserRoleCostumer:
		return entity.UserKindCustomer
	case entity.UserRoleStoreOwner, entity.UserRoleStoreEmployee:
		return entity.UserKindStore
	case entity.UserRoleAdmin:
		return entity.UserKindAdmin
	case entity.UserRoleSupport:
		return entity.UserKindSupport
	default:
		// default seguro: customer
		return entity.UserKindCustomer
	}
}

func decideNextStep(kind entity.UserKind, storesCount int) entity.NextStep {
	if kind == entity.UserKindCustomer {
		return entity.NextStepBrowseStores
	}
	// kind == store
	if storesCount <= 0 {
		return entity.NextStepCreateStore
	}
	return entity.NextStepStoreDashboard
}
