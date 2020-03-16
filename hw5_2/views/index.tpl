<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
</head>

<body>
    <div>
        <h1>{{.Title}}</h1>
        <div>
            <ul>
                {{range .Posts}}
                <li>
                    <div>
                        <h3>{{.Title}}</h3>
                        <h4>{{.Date}}</h4>
                        <p>{{.Content}}</p>
                        <p>{{.Link}}</p>
                        <a href="/post/?id={{.ID}}">Read</a>
                        <a href="/edit/?id={{.ID}}">Edit</a>
                    </div>
                </li>
                {{end}}
            </ul>
        </div>
    </div>
</body>

</html>
