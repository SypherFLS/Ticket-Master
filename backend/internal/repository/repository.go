package repository

import (

)

type Repository interface {
	HealthCheck()
}