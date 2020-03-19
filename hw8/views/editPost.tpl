<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <title>{{.Title}}</title>
</head>

<body>
    <div class="container">
        <h3>Edit Post</h3>
        <form method="POST" action="/edit">
            <table>
                <tr>
                    <td style="display:none;">
                        <label>Id</label>
                        <input type="id" name="id" value="{{.Post.ID}}">
                    </td>
                    <td>
                        <label>Title</label>
                        <input type="title" name="title" value="{{.Post.Title}}">
                    </td>
                    <td>
                        <label>Date</label>
                        <input type="date" name="date" value="{{.Post.Date}}">
                    </td>
                </tr>
                <tr>
                    <td colspan="3">
                        <label>Link</label>
                        <input type="text" name="link" value="{{.Post.Link}}">
                    </td>
                </tr>
                <tr>
                    <td colspan="3">
                        <label>Content</label><br>
                        <textarea name="content" rows="10" cols="40">{{.Post.Content}}</textarea>
                    </td>
                </tr>
            </table>
            <input type="submit" value="submit">
            <a href="/">Back</a>
        </form>
    </div>
</body>

</html>