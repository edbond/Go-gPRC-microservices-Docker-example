module clientservice

go 1.15

require github.com/sirupsen/logrus v1.7.0

require (
	google.golang.org/grpc v1.35.0
	ports.services.com/ports v0.0.1
)

replace ports.services.com/ports => ../ports
