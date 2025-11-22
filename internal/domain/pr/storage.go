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
	// // (POST /pullRequest/merge)
	// PullRequestMerge()
	// // Переназначить конкретного ревьювера на другого из его команды
	// // (POST /pullRequest/reassign)
	// PullRequestReassign()
	// // Создать команду с участниками (создаёт/обновляет пользователей)
	// // (POST /team/add)
	// TeamAdd()
	// // Получить команду с участниками
	// // (GET /team/get)
	// TeamGet()
	// // Получить PR'ы, где пользователь назначен ревьювером
	// // (GET /users/getReview)
	// UsersGetReview()
	// // Установить флаг активности пользователя
	// // (POST /users/setIsActive)
	// UsersSetIsActive()
}
