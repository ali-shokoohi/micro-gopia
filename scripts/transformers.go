package scripts

import (
	"encoding/json"

	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/entities"
)

// UserEntityToUserViewDto transforms an UserEntity to an UserViewDto
func UserEntityToUserViewDto(userEntity *entities.UserEntity) (*dto.UserViewDto, error) {
	user, err := json.Marshal(userEntity)
	if err != nil {
		return nil, err
	}
	var userViewDto *dto.UserViewDto
	err = json.Unmarshal(user, &userViewDto)
	if err != nil {
		return nil, err
	}
	return userViewDto, nil
}

// UserEntitiesToUserViewDtos transforms a slice of UserEntities to a slice of UserViewDtos
func UserEntitiesToUserViewDtos(userEntities []*entities.UserEntity) ([]*dto.UserViewDto, error) {
	var userViewDtos []*dto.UserViewDto
	if userEntities == nil {
		return userViewDtos, nil
	}
	for _, userEntity := range userEntities {
		userViewDto, err := UserEntityToUserViewDto(userEntity)
		if err != nil {
			return nil, err
		}
		userViewDtos = append(userViewDtos, userViewDto)
	}
	return userViewDtos, nil
}

func ErrorsSliceToStringSlice(errs []error) []string {
	var strs []string
	if errs == nil {
		return strs
	}
	for _, err := range errs {
		if err != nil {
			strs = append(strs, err.Error())
		}
	}
	return strs
}
