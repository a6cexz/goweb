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
            <h3>{{.Post.Title}}</h3>
            <h4>{{.Post.Date}}</h4>
            <p>{{.Post.Content}}</p>
            <p>{{.Post.Link}}</p>
            <a href="/edit/?id={{.Post.ID}}">Edit</a>
        </div>
    </div>
</body>

</html>
