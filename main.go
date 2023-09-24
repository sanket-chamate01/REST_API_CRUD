package main

func main() {
	app := App{}
	app.Initialise(DBUser, DBPassword, DBName)
	app.Run("localhost:10000")
}
