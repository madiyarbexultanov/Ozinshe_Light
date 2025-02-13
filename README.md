## Ozinshe

A personal online platform providing movie information. The system allows users to rate movies, mark them as watched, and create a watchlist.

### Functional Requirements

* Create, edit, and delete movies and their details, including title, description, director, release year, genre, trailer link, and poster;
* Sort and filter movies based on various criteria;
* Rate movies;
* Create a watchlist;
* Mark movies as watched;
* Create, edit, and delete genres;
* Create, edit, reset passwords, and delete users;
* Users must log in with an email and password to access the system.

### Нефункциональные требования

* Provide a documented API following the OpenAPI specification;
* Support containerization via Docker.

## How to Run the UI?

To run the Ozinshe UI locally, you need Git and Docker.

### Running the UI Without Authentication

```
docker run --name ozinshe-ui -p "8080:3000" -e VITE_API_URL="http://localhost:8081" -e VITE_FEATURE_AUTH="false" -e VITE_SIMPLIFIED_MOVIE="true" -d kchsherbakov/ozinshe-ui:latest
```

### Running the UI with Authentication Enabled
```
docker run --name ozinshe-ui -p "8080:3000" -e VITE_API_URL="http://localhost:8081" -e VITE_FEATURE_AUTH="true" -e VITE_SIMPLIFIED_MOVIE="false" -d kchsherbakov/ozinshe-ui:latest
```

To log in, use the following credentials:

Email: admin@admin.com
Password: admin