package pr

type Storage interface {
	// Создать PR и автоматически назначить до 2 ревьюверов из команды автора
	PullRequestCreate(pr PullRequest) error
	// Установить флаг активности пользователя
	UsersSetIsActive(id string) error
	//Полуить команду автора
	GetAuthorTeam(id string) (string, error)
	//Получить свободных ревьеров
	GetFreeReviewers(team string, authorid string) ([]User, error)
	// // Пометить PR как MERGED (идемпотентная операция)
	PullRequestMerge(id string) (PullRequest, error)
	// // Переназначить конкретного ревьювера на другого из его команды)
	PullRequestReassign(r PostPullRequestReassign) (PullRequest, error)
	// // Создать команду с участниками (создаёт/обновляет пользователей)
	TeamAdd(t Team) (Team, error)
	// // Получить команду с участниками
	// // (GET /team/get)
	TeamGet(teamName string)(Team, error)
	// // Получить PR'ы, где пользователь назначен ревьювером
	// // (GET /users/getReview)
	// UsersGetReview()
	// // Установить флаг активности пользователя
	// // (POST /users/setIsActive)
	// UsersSetIsActive()
}
