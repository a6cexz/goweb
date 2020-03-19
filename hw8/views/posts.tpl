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