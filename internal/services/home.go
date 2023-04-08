package services

type HomeService struct{}

// GetStartTime returns the current time as a string
func (s *HomeService) Root() string {
	return "hello"
}
