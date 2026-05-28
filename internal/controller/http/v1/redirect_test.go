package v1

import (
	"drk-url-shortener/internal/usecase/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_router_redirect(t *testing.T) {
	// type args struct {
	// 	w   http.ResponseWriter
	// 	req *http.Request
	// }

	// mockBehavior (имитация поведения), тип-функция (function type), настройщик поведения мока.
	// Callback-функция — так как эта логика передается внутрь теста, чтобы сработать в нужный момент (инъекция поведения).
	// В данном случае принимает объект (структуру) имитатора ShortURL interface и строку слага.
	type mockBehavior func(s *mocks.MockShortURL, slug string)

	tests := []struct {
		name string
		// r            router
		// args         args
		id           string
		mockBehavior mockBehavior
	}{
		{
			name: "OK",
			id:   "abc",
			mockBehavior: func(s *mocks.MockShortURL, slug string) {
				// У объекта имитатора s вызываем метод EXPECT (создаем поведение для s)
				// s.EXPECT() активирует режим записи (recorder *MockShortURLMockRecorder — указатель на "записывающий" объект,
				// который предоставляет методы для настройки ожидаемых вызовов (EXPECT().MethodName(arg).Return(value))
				// Сообщаем имитатору: «Сейчас я опишу вызов, который должен произойти во время работы программы».
				s.EXPECT(). // (здесь) при обращении к объекту s мы будем ОЖИДАТЬ()
						GetOriginal(slug).                 // Указывает, какой именно метод мы ждем и с какими аргументами
						Return("https://google.com", nil). // Определяет результат, который метод вернет тестируемому коду.
					// В данном случае,
					// как только программа вызовет метод GetOriginal,
					// имитатор мгновенно отдаст ей "https://google.com" (url) и nil (отсутствие ошибки).
					// Это позволяет тестировать логику дальше, не обращаясь к реальной базе данных.
					Times(1) // (не обязательно) - сколько раз вызываем (по умолчанию 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализация моков.
			ctrl := gomock.NewController(t)
			// Создаем "ложный" сервис, который "притворяется" реальной бизнес-логикой (интерфейсом ShortURL).
			ShortURL := mocks.NewMockShortURL(ctrl)
			// Передавая параметры (ShortURL, tt.id) мы указываем, что
			// 🕖 ожидаем получить вызов методов сервиса ShortURL, а в качестве аргумента передадим id
			// ❗Значит и должен быть вызван только метод принимающий такой параметр!
			tt.mockBehavior(ShortURL, tt.id)

			// // Создаем объект сервисов, но передадим аргументом для интерфейса авторизации наш "ложный" auth.
			// // service.Service "думает", что работает с настоящей базой или API, хотя на самом деле он работает с контролируемым нами моком.
			// services := &service.Service{Authorization: auth}
			// // Инициализируем хендлер.
			// // Структура controller получает объект services, внутри которого уже есть наш мок.
			// // Хендлер не знает, как регистрировать пользователя, он лишь делегирует это сервису.
			// handler := controller{services}

			// // 2. Init Endpoint
			// r := gin.New()
			// // Регистрируем конкретную функцию контроллера (тестируемый метод хендлера signUp) на маршрут /register
			// // Gin теперь знает: «Если придет POST-запрос на этот адрес, нужно запустить именно этот код».
			// // А при вызове handler.signUp внутри сработает цепочка, ведущая к моку.
			// r.POST("/register", handler.signUp)

			// // Создадим запрос к тестируемуму методу
			// w := httptest.NewRecorder()
			// req := httptest.NewRequest("POST", "/register",
			// 	//bytes.NewBufferString(test.inputBody) создает тело нашего запроса
			// 	//и реализует интерфейс io.Reader
			// 	bytes.NewBufferString(test.inputBody))

			// // Make Request
			// // Вызываем метод ServeHTTP у объекта роутера
			// // Отдаем роутеру «виртуальный» запрос (req) и «записывающее устройство» (w).
			// // Роутер прогоняет запрос через свои механизмы, вызывает хендлер, тот вызывает сервис (мок!),
			// // получает ответ и записывает результат в w.
			// r.ServeHTTP(w, req)
		})
	}
}
