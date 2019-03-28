<div align="center">
    <h1>~ discordgocmds ~</h1>
    <strong>Simple command parser package for discordgo</strong><br><br>
    <a href="https://godoc.org/github.com/zekroTJA/discordgocmds"><img src="https://img.shields.io/badge/docs-godoc-c918cc.svg" /></a>
    <!-- <a href="https://www.npmjs.com/package/discordjs-cmds2" ><img src="https://img.shields.io/npm/v/discordjs-cmds2.svg" /></a>&nbsp;
    <a href="https://www.npmjs.com/package/discordjs-cmds2" ><img src="https://img.shields.io/npm/dt/discordjs-cmds2.svg" /></a> -->
<br>
</div>

---

<div align="center">
    <code>go get github.com/zekroTJA/discordgocmds</code>
</div>

---

# [👉 GODOC](https://godoc.org/github.com/zekroTJA/discordgocmds)

---

# Features

- Simply create commands with invoke aliases, permission level, description and help text by creating a class for each command exteding the `Command` abstract class.
- Use whatever you want as Database source to manage permissions and guild prefixes by creating your own database driver extending the `DatabaseInterface` class.
- Permission system using permission levels.
- You can also implement your own permission system into discordjs-cmds2 by using the `PermissionInterface` class.
- Group your commands together
- Automatically created command list and help message
- You can also replace the default help command with your own just by overwriting the `help` invoke.
- Promise-Based safety: Every command will be executed in seperate threads which also will catch all exceptions thrown in the commands code.
- Register your own logger classes based on the `LoggerInterface` if you want to log into a Database or whatever you want to do with it

---

---

© 2018 zekro Development  
[zekro.de](htttps://zekro.de) | contact[at]zekro.de  
MIT Licence