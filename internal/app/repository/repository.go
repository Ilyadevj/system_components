package repository

type ComponentStatus string

const (
	StatusDraft     ComponentStatus = "draft"
	StatusPublished ComponentStatus = "published"
	StatusArchived  ComponentStatus = "archived"
)

// SystemComponent описывает сущность услуги/компонента в 1НФ (без массивов внутри предметных полей)
type SystemComponent struct {
	ID           int
	Title        string
	Description  string
	Priority     int             // Поле по теме 1 (число)
	Grouping     string          // Поле по теме 2 (группировка: Кэш, База данных, Брокер и т.д.)
	Status       ComponentStatus // Статус: draft, published, archived
	ImageURL     string          // Ссылка на MinIO
	VideoURL     string          // Ссылка на MinIO
	LikesUserIDs []int           // Вложенный массив ID пользователей, поставивших лайк
}

type Repository struct {
	items []SystemComponent
}

func NewRepository() *Repository {
	return &Repository{
		items: []SystemComponent{
			{
				ID:           1,
				Title:        "Сервис авторизации и сессий",
				Description:  "Высоконагруженный сервис обработки JWT токенов, аутентификации пользователей и валидации прав доступа в реальном времени с поддержкой ролевой модели.",
				Priority:     1,
				Grouping:     "Кэш",
				Status:       StatusPublished,
				ImageURL:     "http://localhost:9000/media/components/preview-auth.jpg",
				VideoURL:     "http://localhost:9000/media/components/benchmark-auth.mp4",
				LikesUserIDs: []int{101, 102, 105},
			},
			{
				ID:           2,
				Title:        "Кэш-брокер Redis Sentinel",
				Description:  "Распределенный отказоустойчивый кэш оперативной памяти для временных срезов аналитики, снижающий время отклика базы данных до субмиллисекунд.",
				Priority:     2,
				Grouping:     "Кэш",
				Status:       StatusPublished,
				ImageURL:     "http://localhost:9000/media/components/preview-redis.jpg",
				VideoURL:     "http://localhost:9000/media/components/benchmark-redis.mp4",
				LikesUserIDs: []int{101, 104},
			},
			{
				ID:           3,
				Title:        "Кластер PostgreSQL Primary",
				Description:  "Реляционное хранилище транзакционных данных с физической репликацией, автоматическим переключением при сбоях и пулом соединений PgBouncer.",
				Priority:     3,
				Grouping:     "База данных",
				Status:       StatusPublished,
				ImageURL:     "http://localhost:9000/media/components/preview-postgres.jpg",
				VideoURL:     "http://localhost:9000/media/components/benchmark-postgres.mp4",
				LikesUserIDs: []int{102, 103, 104, 105},
			},
			{
				ID:           4,
				Title:        "Устаревший шлюз биллинга",
				Description:  "Архивный компонент обработки платежей через старый протокол, выведенный из продуктовой эксплуатации.",
				Priority:     4,
				Grouping:     "База данных",
				Status:       StatusArchived, // Удален — никогда не выводится в интерфейсе
				ImageURL:     "http://localhost:9000/media/components/preview-postgres.jpg",
				VideoURL:     "http://localhost:9000/media/components/benchmark-postgres.mp4",
				LikesUserIDs: []int{},
			},
			{
				ID:           5,
				Title:        "Аналитическая СУБД ClickHouse",
				Description:  "Колоночная база данных для сбора логов телеметрии и высокоскоростного построения дашбордов активности компонентов инфраструктуры.",
				Priority:     5,
				Grouping:     "База данных",
				Status:       StatusDraft, // Единственный черновик для страницы добавления
				ImageURL:     "http://localhost:9000/media/components/preview-react.jpg",
				VideoURL:     "http://localhost:9000/media/components/benchmark-react.mp4",
				LikesUserIDs: []int{},
			},
		},
	}
}

// GetAvailableGroups возвращает список всех уникальных групп среди опубликованных компонентов
func (r *Repository) GetAvailableGroups() []string {
	groupMap := make(map[string]bool)
	var groups []string

	for _, item := range r.items {
		if item.Status == StatusPublished && item.Grouping != "" {
			if !groupMap[item.Grouping] {
				groupMap[item.Grouping] = true
				groups = append(groups, item.Grouping)
			}
		}
	}
	return groups
}

// GetPublished фильтрует опубликованные компоненты по выбранным галочкам
func (r *Repository) GetPublished(selectedGroups []string) []SystemComponent {
	var result []SystemComponent
	filterMap := make(map[string]bool)
	for _, g := range selectedGroups {
		filterMap[g] = true
	}

	for _, item := range r.items {
		if item.Status != StatusPublished {
			continue // Исключаем draft и archived
		}
		if len(selectedGroups) == 0 || filterMap[item.Grouping] {
			result = append(result, item)
		}
	}
	return result
}

// GetDraft возвращает единственный черновик
func (r *Repository) GetDraft() *SystemComponent {
	for _, item := range r.items {
		if item.Status == StatusDraft {
			return &item
		}
	}
	return nil
}

// GetFeedItem возвращает компонент по ID или циклически переходит к следующему (?next=true)
func (r *Repository) GetFeedItem(targetID int, getNext bool) SystemComponent {
	var published []SystemComponent
	for _, item := range r.items {
		if item.Status == StatusPublished {
			published = append(published, item)
		}
	}
	if len(published) == 0 {
		return SystemComponent{}
	}

	if targetID <= 0 {
		return published[0]
	}

	for i, item := range published {
		if item.ID == targetID {
			if getNext {
				return published[(i+1)%len(published)]
			}
			return item
		}
	}
	return published[0]
}