# Yuri Web Frontend Cookie Usage

Yuri is using permanent cookies to atuomatically authenticate your identity against the back end API and for saving user specific preferences like the way of sorting the sounds list or notification preferences.

On logging in and authorizing access to the Discord App of Yuri, two cookies will be saved.  
- The first one with the name `token` saves the raw automatically generated user token. **Attention:** this value should be treated very securely because this token can be used to authenticate as **you** against Yuri's API in combination with your Discord user ID.  
- The second one named `user_id` contains your Discord user ID.  
Both cookies will be send to the server to authenticate you for following REST API requests and for initializing the web socket connection.

The following settings tokens will be saved:
- `sort_by` - Contains `NAME` or `DATE` depending on which you have set recently with the `SORT BY ...` button.
- `cookies_accepted` - Contains `1` if you have accepted the cookie information notification so that the window will not be shown next time.

Yuri **does not** set cookies to track you or to analyze any of your behaviours.

You can delete all set cookies using the following static page *(requires JavaScript to be enabled)*:  
```
<host>/static/delete-cookies.html
```
