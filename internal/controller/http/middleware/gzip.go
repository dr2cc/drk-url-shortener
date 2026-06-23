package mw

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipReadCloser объединяет gzip.Reader и оригинальное тело запроса
type gzipReadCloser struct {
	gz *gzip.Reader
	rc io.ReadCloser
}

func (g *gzipReadCloser) Read(p []byte) (n int, err error) {
	return g.gz.Read(p)
}

func (g *gzipReadCloser) Close() error {
	// Сначала закрываем gzip reader, чтобы освободить его ресурсы
	if err := g.gz.Close(); err != nil {
		g.rc.Close() // даже при ошибке, пытаемся закрыть основное тело
		return err
	}
	// Затем закрываем оригинальное тело запроса
	return g.rc.Close()
}

func DecompressRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие gzip в заголовке (учитываем регистр и возможные цепочки)
		if strings.Contains(strings.ToLower(r.Header.Get("Content-Encoding")), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress gzip body", http.StatusBadRequest)
				return
			}

			// Подменяем r.Body нашей оберткой.
			// Теперь хендлер сам закроет всё цепочкой в конце работы.
			r.Body = &gzipReadCloser{gz: gz, rc: r.Body}

			// Удаляем заголовки, так как данные теперь распакованы
			r.Header.Del("Content-Encoding")
			r.Header.Del("Content-Length")
		}

		next.ServeHTTP(w, r)
	})
}
