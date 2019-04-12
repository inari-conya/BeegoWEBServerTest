package test_mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	"test"

	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
)

type T_user struct{}

type T_student struct{}

var conn string

func init() {
	test.UserReg("mysql", &T_user{})
	test.StudentReg("mysql", &T_student{})

	b, err := ioutil.ReadFile("server.config")
	if err != nil {
		fmt.Println("打开配置文件失败")
		panic(err)
	}

	js, err := simplejson.NewJson(b)
	if err != nil {
		panic(err.Error())
	}

	ip, _ := js.Get("db_config").Get("ip").String()
	port, _ := js.Get("db_config").Get("port").String()
	dbhost := ip + ":" + port
	dbuser, _ := js.Get("db_config").Get("username").String()
	dbpassword, _ := js.Get("db_config").Get("passwd").String()
	dbmotion_config, _ := js.Get("db_config").Get("databasename").String()

	conn = dbuser + ":" + dbpassword + "@tcp(" + dbhost + ")/" + dbmotion_config + "?charset=utf8"

}

/*打开数据库*/
func db_open() *sql.DB {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err)
		fmt.Println("Open mysql database failed")
	}
	return db
}

func (s *T_student) ListTutorial() ([]map[string]string, error) {
	db := db_open()
	defer db.Close()
	l := make([]map[string]string, 0)

	rows, err := db.Query("select fname,lname,age from test.student where tutorial=1;")
	if err != nil {
		fmt.Println("query db failed")
		return l, err
	}
	defer rows.Close()

	col, err := rows.Columns()
	if err != nil {
		fmt.Println("get columns failed")
		return l, err
	}

	cols := make([]string, len(col))
	for k, v := range col {
		cols[k] = strings.ToLower(v)
	}

	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	for rows.Next() {
		rows.Scan(scans...)

		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			vv := string(v)
			row[key] = vv
		}
		l = append(l, row)
	}

	return l, nil
}

func (u *T_user) Auth(id, authstr string) bool {
	db := db_open()
	defer db.Close()

	rows, err := db.Query("select id from test.int_auth_token_cache where id=? and int_auth_token=?;", id, authstr)
	if err != nil {
		fmt.Println("Auth query failed", err)
		return false
	}
	defer rows.Close()

	return rows.Next()
}

func (u *T_user) Detail(id, authstr string) (map[string]string, error) {
	m := make(map[string]string)

	db := db_open()
	defer db.Close()

	rows, err := db.Query("select user.id,user.name,int_auth_token_cache.ip,int_auth_token_cache.device from test.int_auth_token_cache inner join test.user on int_auth_token_cache.id=user.id where user.id=? and int_auth_token_cache.int_auth_token=?;", id, authstr)
	if err != nil {
		fmt.Println("Detail query failed", err)
		return m, err
	}
	defer rows.Close()

	col, err := rows.Columns()
	if err != nil {
		fmt.Println("get columns failed")
		return m, err
	}

	cols := make([]string, len(col))
	for k, v := range col {
		cols[k] = strings.ToLower(v)
	}

	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	for rows.Next() {
		rows.Scan(scans...)

		for k, v := range vals {
			key := cols[k]
			vv := string(v)
			m[key] = vv
		}
	}
	return m, nil
}

// /*返回用户密码*/
// func (u *T_user) Getpswd(username string) (pswd string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select password from t_user where username=?", username)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer rows.Close()
// 	ok := rows.Next()
// 	if ok {
// 		rows.Scan(&pswd)
// 	} else {
// 		err = errors.New("user name does not exist")
// 		return "", err
// 	}
// 	return pswd, nil
// }

// //返回用户名称
// func (u *T_user) Getname(username string) (name string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select nickname from t_user where username=?", username)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer rows.Close()
// 	ok := rows.Next()
// 	if ok {
// 		rows.Scan(&name)
// 	} else {
// 		err = errors.New("user name does not exist")
// 		return "", err
// 	}
// 	return name, nil
// }

// /*判断用户是否已存在*/
// func (u *T_user) Exist(username string) bool {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select password from t_user where username=?", username)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return false
// 	}
// 	defer rows.Close()
// 	return rows.Next()
// }

// /*展示所有已存在的用户*/
// func (u *T_user) Showallusers() (users, names []string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select username,nickname from t_user")
// 	if err != nil {
// 		return users, names, err
// 	}
// 	defer rows.Close()
// 	var usertmp, nametmp string
// 	for rows.Next() {
// 		err = rows.Scan(&usertmp, &nametmp)
// 		if err != nil {
// 			return users, names, err
// 		}
// 		users = append(users, usertmp)
// 		names = append(names, nametmp)
// 	}
// 	return users, names, nil
// }

// /*更新用户名*/
// func (u *T_user) Updatename(username, name string) error {
// 	db := db_open()
// 	defer db.Close()

// 	if name == "" {
// 		return errors.New("nickname can not be empty")
// 	}

// 	r, err := db.Exec("update t_user set nickname = ? where username = ?", name, username)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Updatename RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("userid does not exist")
// 	}
// 	return nil
// }

// /*更新用户密码*/
// func (u *T_user) Updatepswd(username, pswd string) error {
// 	db := db_open()
// 	defer db.Close()

// 	if pswd == "" {
// 		return errors.New("password can not be empty")
// 	}

// 	r, err := db.Exec("update t_user set password = ? where username = ?", pswd, username)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Updatepswd RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("userid does not exist")
// 	}
// 	return nil
// }

// /*插入新用户*/
// func (u *T_user) Insert(username, pswd, name string) error {
// 	db := db_open()
// 	defer db.Close()

// 	if pswd == "" || name == "" {
// 		return errors.New("password or nickname can not be empty")
// 	}

// 	r, err := db.Exec(`insert into t_user (username,password,nickname) values (?,?,?)`, username, pswd, name)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Insertuser RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("userid already exists")
// 	}
// 	return nil
// }

// /*删除用户*/
// func (u *T_user) Delete(username string) error {
// 	db := db_open()
// 	defer db.Close()

// 	r, err := db.Exec("delete from t_user where username=?", username)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Delete user RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("userid does not exist")
// 	}

// 	_, err = db.Exec("delete from user_role where userid=?", username)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec("delete from group_user where userid=?", username)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// /*展示所有已存在的角色*/
// func (u *T_user) Showallroles() (ids []int, names, descs []string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id,name,description from t_role")
// 	if err != nil {
// 		return ids, names, descs, err
// 	}
// 	defer rows.Close()
// 	var nametmp, desctmp string
// 	var idtmp int
// 	for rows.Next() {
// 		err = rows.Scan(&idtmp, &nametmp, &desctmp)
// 		if err != nil {
// 			return ids, names, descs, err
// 		}
// 		ids = append(ids, idtmp)
// 		names = append(names, nametmp)
// 		descs = append(descs, desctmp)
// 	}
// 	return ids, names, descs, nil
// }

// /*判断角色是否已存在*/
// func (r *T_role) Exist(name string) bool {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id from t_role where name=?", name)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return false
// 	}
// 	defer rows.Close()
// 	return rows.Next()
// }

// /*插入新角色*/
// func (u *T_role) Insert(name, desc string) error {
// 	db := db_open()
// 	defer db.Close()

// 	if u.Exist(name) {
// 		return errors.New("role name already exists")
// 	}

// 	r, err := db.Exec(`insert into t_role (name,description) values (?,?)`, name, desc)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Insertuser RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("role name already exists")
// 	}
// 	return nil
// }

// /*更新角色描述*/
// func (u *T_role) Updatedesc(name, desc string) error {
// 	db := db_open()
// 	defer db.Close()

// 	r, err := db.Exec("update t_role set description = ? where name = ?", desc, name)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Updatename RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("role name does not exist or description has not changed")
// 	}
// 	return nil
// }

// /*删除角色*/
// func (u *T_role) Delete(name string) error {
// 	db := db_open()
// 	defer db.Close()

// 	id, err := u.GetID(name)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = db.Exec("delete from user_role where roleid=?", id)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec("delete from role_privilege where roleid=?", id)
// 	if err != nil {
// 		return err
// 	}

// 	r, err := db.Exec("delete from t_role where id=?", id)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Delete user RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("role name does not exist")
// 	}
// 	return nil
// }

// /*获取角色ID*/
// func (u *T_role) GetID(name string) (id int, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id from t_role where name=?", name)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return 0, err
// 	}
// 	defer rows.Close()
// 	ok := rows.Next()
// 	if ok {
// 		rows.Scan(&id)
// 	} else {
// 		err = errors.New("role name does not exist")
// 		return 0, err
// 	}
// 	return id, nil
// }

// /*展示所有已存在的权限*/
// func (u *T_user) Showallprvgs() (ids []int, names, descs []string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id,name,description from t_privilege")
// 	if err != nil {
// 		return ids, names, descs, err
// 	}
// 	defer rows.Close()
// 	var nametmp, desctmp string
// 	var idtmp int
// 	for rows.Next() {
// 		err = rows.Scan(&idtmp, &nametmp, &desctmp)
// 		if err != nil {
// 			return ids, names, descs, err
// 		}
// 		ids = append(ids, idtmp)
// 		names = append(names, nametmp)
// 		descs = append(descs, desctmp)
// 	}
// 	return ids, names, descs, nil
// }

// /*判断权限是否已存在*/
// func (r *T_privilege) Exist(name string) bool {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id from t_privilege where name=?", name)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return false
// 	}
// 	defer rows.Close()
// 	return rows.Next()
// }

// /*插入新权限*/
// func (u *T_privilege) Insert(name, desc string) error {
// 	db := db_open()
// 	defer db.Close()

// 	if u.Exist(name) {
// 		return errors.New("privilege name already exists")
// 	}

// 	r, err := db.Exec(`insert into t_privilege (name,description) values (?,?)`, name, desc)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Insertuser RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("privilege name already exists")
// 	}
// 	return nil
// }

// /*更新权限描述*/
// func (u *T_privilege) Updatedesc(name, desc string) error {
// 	db := db_open()
// 	defer db.Close()

// 	r, err := db.Exec("update t_privilege set description = ? where name = ?", desc, name)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Updatename RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("privilege name does not exist or description has not changed")
// 	}
// 	return nil
// }

// /*删除权限*/
// func (u *T_privilege) Delete(name string) error {
// 	db := db_open()
// 	defer db.Close()

// 	id, err := u.GetID(name)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = db.Exec("delete from role_privilege where prvgid=?", id)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec("delete from group_privilege where prvgid=?", id)
// 	if err != nil {
// 		return err
// 	}

// 	r, err := db.Exec("delete from t_privilege where id=?", id)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Delete user RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("privilege name does not exist")
// 	}
// 	return nil
// }

// /*获取权限ID*/
// func (u *T_privilege) GetID(name string) (id int, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id from t_privilege where name=?", name)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return 0, err
// 	}
// 	defer rows.Close()
// 	ok := rows.Next()
// 	if ok {
// 		rows.Scan(&id)
// 	} else {
// 		err = errors.New("privilege name does not exist")
// 		return 0, err
// 	}
// 	return id, nil
// }

// /*展示所有已存在的组*/
// func (u *T_user) Showallgroups() (ids []int, names, descs []string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id,name,description from t_group")
// 	if err != nil {
// 		return ids, names, descs, err
// 	}
// 	defer rows.Close()
// 	var nametmp, desctmp string
// 	var idtmp int
// 	for rows.Next() {
// 		err = rows.Scan(&idtmp, &nametmp, &desctmp)
// 		if err != nil {
// 			return ids, names, descs, err
// 		}
// 		ids = append(ids, idtmp)
// 		names = append(names, nametmp)
// 		descs = append(descs, desctmp)
// 	}
// 	return ids, names, descs, nil
// }

// /*判断组是否已存在*/
// func (r *T_group) Exist(name string) bool {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id from t_group where name=?", name)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return false
// 	}
// 	defer rows.Close()
// 	return rows.Next()
// }

// /*插入新组*/
// func (u *T_group) Insert(name, desc string) error {
// 	db := db_open()
// 	defer db.Close()

// 	if u.Exist(name) {
// 		return errors.New("group name already exists")
// 	}

// 	r, err := db.Exec(`insert into t_group (name,description) values (?,?)`, name, desc)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Insertuser RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("group name already exists")
// 	}
// 	return nil
// }

// /*更新组描述*/
// func (u *T_group) Updatedesc(name, desc string) error {
// 	db := db_open()
// 	defer db.Close()

// 	r, err := db.Exec("update t_group set description = ? where name = ?", desc, name)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Updatename RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("group name does not exist or description has not changed")
// 	}
// 	return nil
// }

// /*删除组*/
// func (u *T_group) Delete(name string) error {
// 	db := db_open()
// 	defer db.Close()

// 	id, err := u.GetID(name)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = db.Exec("delete from group_user where groupid=?", id)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec("delete from group_privilege where groupid=?", id)
// 	if err != nil {
// 		return err
// 	}

// 	r, err := db.Exec("delete from t_group where id=?", id)
// 	if err != nil {
// 		return err
// 	}
// 	n, err1 := r.RowsAffected()
// 	if err1 != nil {
// 		fmt.Println("Delete user RowsAffected error:", err)
// 	}
// 	if n == 0 {
// 		return errors.New("group name does not exist")
// 	}
// 	return nil
// }

// /*获取组ID*/
// func (u *T_group) GetID(name string) (id int, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select id from t_group where name=?", name)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return 0, err
// 	}
// 	defer rows.Close()
// 	ok := rows.Next()
// 	if ok {
// 		rows.Scan(&id)
// 	} else {
// 		err = errors.New("group name does not exist")
// 		return 0, err
// 	}
// 	return id, nil
// }

// /*展示所有用户角色关联*/
// func (u *T_user) Showallurs() (map[string][]string, error) {
// 	urs := make(map[string][]string)
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select userid from user_role")
// 	if err != nil {
// 		return urs, err
// 	}
// 	defer rows.Close()
// 	var uid, rname string
// 	var uids, rnames []string
// 	for rows.Next() {
// 		err = rows.Scan(&uid)
// 		if err != nil {
// 			return urs, err
// 		}
// 		uids = append(uids, uid)
// 	}
// 	for _, v := range uids {
// 		rows1, err1 := db.Query("select roleid from user_role where userid=?", v)
// 		if err1 != nil {
// 			return urs, err1
// 		}
// 		defer rows1.Close()
// 		var rid int
// 		for rows1.Next() {
// 			err = rows1.Scan(&rid)
// 			if err != nil {
// 				return urs, err
// 			}
// 			rows2, err2 := db.Query("select name from t_role where id=?", rid)
// 			if err2 != nil {
// 				return urs, err2
// 			}
// 			defer rows2.Close()
// 			for rows2.Next() {
// 				err = rows2.Scan(&rname)
// 				if err != nil {
// 					return urs, err
// 				}
// 				rnames = append(rnames, rname)
// 			}
// 		}
// 		urs[v] = rnames
// 		rnames = make([]string, 0)
// 	}
// 	return urs, nil
// }

// /*插入新用户角色关联*/
// func (u *T_user) Setroles(username string, ids []int) error {
// 	db := db_open()
// 	defer db.Close()

// 	for _, rid := range ids {
// 		r, err := db.Exec(`insert into user_role (userid,roleid) values (?,?)`, username, rid)
// 		if err != nil {
// 			return err
// 		}
// 		n, err1 := r.RowsAffected()
// 		if err1 != nil {
// 			fmt.Println("Insertuser RowsAffected error:", err)
// 		}
// 		if n == 0 {
// 			return errors.New("a relationship betwee user and role already exists")
// 		}

// 	}
// 	return nil
// }

// /*移除新用户角色关联*/
// func (u *T_user) Rmroles(username string, ids []int) error {

// 	db := db_open()
// 	defer db.Close()

// 	for _, rid := range ids {
// 		_, err := db.Exec("delete from user_role where userid=? and roleid=?", username, rid)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// /*展示用户所有关联的角色*/
// func (u *T_user) Getroles(username string) (names []string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select roleid from user_role where userid=?", username)
// 	if err != nil {
// 		return names, err
// 	}
// 	defer rows.Close()
// 	var rid, rname string
// 	for rows.Next() {
// 		err = rows.Scan(&rid)
// 		if err != nil {
// 			return names, err
// 		}
// 		rows1, err1 := db.Query("select name from t_role where id=?", rid)
// 		if err1 != nil {
// 			return names, err1
// 		}
// 		defer rows1.Close()
// 		if rows1.Next() {
// 			err = rows1.Scan(&rname)
// 			if err != nil {
// 				return names, err
// 			}
// 			names = append(names, rname)
// 		}
// 	}
// 	return names, nil
// }

// /*展示所有用户组与用户的关联*/
// func (u *T_user) Showallgus() (map[string][]string, error) {
// 	gus := make(map[string][]string)
// 	t_group := make(map[int]string)
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select groupid from group_user")
// 	if err != nil {
// 		return gus, err
// 	}
// 	defer rows.Close()
// 	var gname, uid string
// 	var uids []string
// 	var gid int
// 	for rows.Next() {
// 		err = rows.Scan(&gid)
// 		if err != nil {
// 			return gus, err
// 		}
// 		rows1, err1 := db.Query("select name from t_group where id=?", gid)
// 		if err1 != nil {
// 			return gus, err1
// 		}
// 		defer rows1.Close()
// 		if rows1.Next() {
// 			err = rows1.Scan(&gname)
// 			if err != nil {
// 				return gus, err
// 			}
// 			t_group[gid] = gname
// 		}
// 	}
// 	for k, v := range t_group {
// 		rows2, err2 := db.Query("select userid from group_user where groupid=?", k)
// 		if err2 != nil {
// 			return gus, err2
// 		}
// 		defer rows2.Close()
// 		for rows2.Next() {
// 			err = rows2.Scan(&uid)
// 			if err != nil {
// 				return gus, err
// 			}
// 			uids = append(uids, uid)
// 		}
// 		gus[v] = uids
// 		uids = make([]string, 0)
// 	}
// 	return gus, nil
// }

// /*插入新用户组-用户关联*/
// func (u *T_group) Setusers(gid int, ids []string) error {
// 	db := db_open()
// 	defer db.Close()

// 	for _, uid := range ids {

// 		r, err := db.Exec(`insert into group_user (groupid,userid) values (?,?)`, gid, uid)
// 		if err != nil {
// 			return err
// 		}
// 		n, err1 := r.RowsAffected()
// 		if err1 != nil {
// 			fmt.Println("Insertuser RowsAffected error:", err)
// 		}
// 		if n == 0 {
// 			return errors.New("a relationship betwee user and role already exists")
// 		}

// 	}
// 	return nil
// }

// /*移除用户组与用户间的关联*/
// func (u *T_group) Rmusers(gid int, ids []string) error {
// 	db := db_open()
// 	defer db.Close()

// 	for _, v := range ids {
// 		_, err := db.Exec("delete from group_user where userid=? and groupid=?", v, gid)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// /*展示用户组下关联的用户*/
// func (u *T_group) Getusers(gid int) (ids, names []string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select userid from group_user where groupid=?", gid)
// 	if err != nil {
// 		return ids, names, err
// 	}
// 	defer rows.Close()
// 	var uid, uname string
// 	for rows.Next() {
// 		err = rows.Scan(&uid)
// 		if err != nil {
// 			return ids, names, err
// 		}
// 		ids = append(ids, uid)
// 		rows1, err1 := db.Query("select nickname from t_user where username=?", uid)
// 		if err1 != nil {
// 			return ids, names, err1
// 		}
// 		defer rows1.Close()
// 		if rows1.Next() {
// 			err = rows1.Scan(&uname)
// 			if err != nil {
// 				return ids, names, err
// 			}
// 			names = append(names, uname)
// 		}
// 	}
// 	return ids, names, nil
// }

// /*展示所有角色与权限的关联*/
// func (u *T_user) Showallrps() (map[string][]string, error) {
// 	rps := make(map[string][]string)
// 	t_role := make(map[int]string)
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select roleid from role_privilege")
// 	if err != nil {
// 		return rps, err
// 	}
// 	defer rows.Close()
// 	var rid int
// 	for rows.Next() {
// 		err = rows.Scan(&rid)
// 		if err != nil {
// 			return rps, err
// 		}
// 		rows1, err1 := db.Query("select name from t_role where id=?", rid)
// 		if err1 != nil {
// 			return rps, err1
// 		}
// 		var rname string
// 		defer rows1.Close()
// 		if rows1.Next() {
// 			err = rows1.Scan(&rname)
// 			if err != nil {
// 				return rps, err
// 			}
// 			t_role[rid] = rname
// 		}
// 	}
// 	for k, v := range t_role {
// 		rows2, err2 := db.Query("select prvgid from role_privilege where roleid=?", k)
// 		if err2 != nil {
// 			return rps, err2
// 		}
// 		defer rows2.Close()
// 		var pid int
// 		for rows2.Next() {
// 			err = rows2.Scan(&pid)
// 			if err != nil {
// 				return rps, err
// 			}
// 			rows3, err3 := db.Query("select name from t_privilege where id=?", pid)
// 			if err3 != nil {
// 				return rps, err3
// 			}
// 			defer rows3.Close()
// 			var pname string
// 			if rows3.Next() {
// 				err = rows3.Scan(&pname)
// 				if err != nil {
// 					return rps, err
// 				}
// 				rps[v] = append(rps[v], pname)
// 			}
// 		}
// 	}
// 	return rps, nil
// }

// /*插入新角色-权限关联*/
// func (u *T_role) Setprvgs(rid int, pids []int) error {
// 	db := db_open()
// 	defer db.Close()

// 	for _, key := range pids {

// 		r, err := db.Exec(`insert into role_privilege (roleid,prvgid) values (?,?)`, rid, key)
// 		if err != nil {
// 			return err
// 		}
// 		n, err1 := r.RowsAffected()
// 		if err1 != nil {
// 			fmt.Println("Insertuser RowsAffected error:", err)
// 		}
// 		if n == 0 {
// 			return errors.New("a relationship betwee role and privilege already exists")
// 		}

// 	}
// 	return nil
// }

// /*移除角色与权限间的关联*/
// func (u *T_role) Rmprvgs(rid int, pids []int) error {
// 	db := db_open()
// 	defer db.Close()

// 	for _, pid := range pids {
// 		_, err := db.Exec("delete from role_privilege where roleid=? and prvgid=?", rid, pid)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// /*展示角色拥有的所有权限*/
// func (u *T_role) Getprvgs(rid int) (names, desc []string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select prvgid from role_privilege where roleid=?", rid)
// 	if err != nil {
// 		return names, desc, err
// 	}
// 	defer rows.Close()
// 	var pname, dsc string
// 	var pid int
// 	for rows.Next() {
// 		err = rows.Scan(&pid)
// 		if err != nil {
// 			return names, desc, err
// 		}
// 		rows1, err1 := db.Query("select name,description from t_privilege where id=?", pid)
// 		if err1 != nil {
// 			return names, desc, err
// 		}
// 		defer rows1.Close()
// 		if rows1.Next() {
// 			err = rows1.Scan(&pname, &dsc)
// 			if err != nil {
// 				return names, desc, err
// 			}
// 			names = append(names, pname)
// 			desc = append(desc, dsc)
// 		}
// 	}
// 	return names, desc, err
// }

// /*展示所有用户组与权限的关联*/
// func (u *T_user) Showallgps() (map[string][]string, error) {
// 	gps := make(map[string][]string)
// 	t_group := make(map[int]string)
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select groupid from group_privilege")
// 	if err != nil {
// 		return gps, err
// 	}
// 	defer rows.Close()
// 	var gid int
// 	for rows.Next() {
// 		err = rows.Scan(&gid)
// 		if err != nil {
// 			return gps, err
// 		}
// 		rows1, err1 := db.Query("select name from t_group where id=?", gid)
// 		if err1 != nil {
// 			return gps, err1
// 		}
// 		var gname string
// 		defer rows1.Close()
// 		if rows1.Next() {
// 			err = rows1.Scan(&gname)
// 			if err != nil {
// 				return gps, err
// 			}
// 			t_group[gid] = gname
// 		}
// 	}
// 	for k, v := range t_group {
// 		rows2, err2 := db.Query("select prvgid from group_privilege where groupid=?", k)
// 		if err2 != nil {
// 			return gps, err2
// 		}
// 		defer rows2.Close()
// 		var pid int
// 		for rows2.Next() {
// 			err = rows2.Scan(&pid)
// 			if err != nil {
// 				return gps, err
// 			}
// 			rows3, err3 := db.Query("select name from t_privilege where id=?", pid)
// 			if err3 != nil {
// 				return gps, err3
// 			}
// 			defer rows3.Close()
// 			var pname string
// 			if rows3.Next() {
// 				err = rows3.Scan(&pname)
// 				if err != nil {
// 					return gps, err
// 				}
// 				gps[v] = append(gps[v], pname)
// 			}
// 		}
// 	}
// 	return gps, nil
// }

// /*插入新用户组-权限关联*/
// func (u *T_group) Setprvgs(gid int, pids []int) error {
// 	db := db_open()
// 	defer db.Close()

// 	for _, key := range pids {

// 		r, err := db.Exec(`insert into group_privilege (groupid,prvgid) values (?,?)`, gid, key)
// 		if err != nil {
// 			return err
// 		}
// 		n, err1 := r.RowsAffected()
// 		if err1 != nil {
// 			fmt.Println("Insertuser RowsAffected error:", err)
// 		}
// 		if n == 0 {
// 			return errors.New("a relationship betwee group and privilege already exists")
// 		}

// 	}
// 	return nil
// }

// /*移除用户组与权限间的关联*/
// func (u *T_group) Rmprvgs(gid int, pids []int) error {
// 	db := db_open()
// 	defer db.Close()

// 	for _, pid := range pids {
// 		_, err := db.Exec("delete from group_privilege where groupid=? and prvgid=?", gid, pid)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// /*展示角色拥有的所有权限*/
// func (u *T_group) Getprvgs(gid int) (names, desc []string, err error) {
// 	db := db_open()
// 	defer db.Close()

// 	rows, err := db.Query("select prvgid from group_privilege where groupid=?", gid)
// 	if err != nil {
// 		return names, desc, err
// 	}
// 	defer rows.Close()
// 	var pname, dsc string
// 	var pid int
// 	for rows.Next() {
// 		err = rows.Scan(&pid)
// 		if err != nil {
// 			return names, desc, err
// 		}
// 		rows1, err1 := db.Query("select name,description from t_privilege where id=?", pid)
// 		if err1 != nil {
// 			return names, desc, err
// 		}
// 		defer rows1.Close()
// 		if rows1.Next() {
// 			err = rows1.Scan(&pname, &dsc)
// 			if err != nil {
// 				return names, desc, err
// 			}
// 			names = append(names, pname)
// 			desc = append(desc, dsc)
// 		}
// 	}
// 	return names, desc, err
// }
