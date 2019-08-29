package storage

import "../shared"

//IRepository - Репозиторий
type IRepository interface
{
	Post(hit shared.GoogleHit) error
	GetCount() int32
}