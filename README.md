#run
go run main.go

#register
http://localhost:9091/api/auth/register

#login
http://localhost:9091/api/auth/login

#paginate
http://localhost:9091/api/v1/elemes

#input form-data file/image upload to cloudinary
http://localhost:9091/api/v1/elemes

#edit
http://localhost:9091/api/v1/elemes/1

#get detail
http://localhost:9091/api/v1/elemes/4

#softdelete
http://localhost:9091/api/v1/elemes/1

#search name AD
http://localhost:9091/api/v1/elemes?limit=10&page=0&sort=created_at desc&name_course.contains=AD

#search in name AD
http://localhost:9091/api/v1/elemes?limit=10&page=0&sort=created_at desc&name_course.in=AD

#search in name AD
http://localhost:9091/api/v1/elemes?limit=10&page=0&sort=created_at desc&name_course.equals=AD
