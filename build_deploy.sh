SET GOOS=linux
go build -o bbljfooddiary.exe
npx warp deploy --project BBLJ --env bbljfooddiary