# openslides-search-service

The OpenSlides search service.

## Configuration:


| Env variable                    | Default value              | Meaning |
| ------------------------------- | -------------------------- | ------- |
| `SECRETS_PATH`                  | `/run/secrets`             | Path where the screts are stored. |
| `OPENSLIDES_SEARCH_PORT`        | `9050`                     | Port the service listens on.    |
| `OPENSLIDES_SEARCH_HOST`        | ``                         | Host the service is bound to.   |
| `OPENSLIDES_SEARCH_MAX_QUEUED`  | `5`                        | Number of waiting queries.      |
| `OPENSLIDES_SEARCH_INDEX_AGE`   | `100ms`                    | Accepted age of internal index. |
| `OPENSLIDES_SEARCH_INDEX_FILE`  | `search.bleve`             | Filename of the internal index. |
| `OPENSLIDES_SEARCH_INDEX_BATCH` | `4096`                     | Batch size of the index when its build or re-generated. |
| `OPENSLIDES_SEARCH_INDEX_UPDATE_INTERVAL` | `120s`           | Poll intervall to update the index without queries. |
| `OPENSLIDES_MODELS_YML`         | `models.yml`               | File path of the used models. |
| `OPENSLIDES_SEARCH_YML`         | `search.yml`               | Fields of the models to be searched. |
| `OPENSLIDES_DB`                 | `openslides`               | Name of the database. |
| `OPENSLIDES_DB_USER`            | `openslides`               | Database user. |
| `OPENSLIDES_DB_PASSWORD`        | `secret:postgres_password` | Password of the database user. |
| `OPENSLIDES_DB_HOST`            | `localhost`                | Host of the database. |
| `OPENSLIDES_DB_PORT`            | `5432`                     | Port of the database. |
| `OPENSLIDES_RESTRICTER`         | ``                         | URL to use the restricter from the auto-update-service to filter the query results.|
