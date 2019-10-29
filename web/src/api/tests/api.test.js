import axios from "axios";

//Axios Config
const withToken = function(config) {
    config.headers.Authorization = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJDck10ZmN1cjFDcVNtT28teHZacUt0ZTRoODk4ZjZpYl9KOGk5TXZDck5zIn0.eyJqdGkiOiIyZDJhZjRjZC1hNWIzLTQ0NzItYTJlYy1jMjJiNjdhZDk5NzEiLCJleHAiOjE1NzIzMTUwODMsIm5iZiI6MCwiaWF0IjoxNTcyMzE0NzgzLCJpc3MiOiJodHRwczovL2lkLnBhdGh3YXIubGFuZC9hdXRoL3JlYWxtcy9QYXRod2FyLURldiIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIyNWFhZGU0My1mOWJlLTQ1YTEtYjI0OC0yZDYwNjU5YzA1MmUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJwbGF0Zm9ybS1mcm9udCIsIm5vbmNlIjoiN2IxYzIxYWItYjJhMi00OWVmLTg3OTItMjFjOGFlNDUxNTdjIiwiYXV0aF90aW1lIjoxNTcyMzE0NzgwLCJzZXNzaW9uX3N0YXRlIjoiNjViODE3MTUtZjVjNC00Yjc4LTg4MTktM2IxNDNkOGM3ODBjIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiQWxvaGEgVGVzdGluZyIsInByZWZlcnJlZF91c2VybmFtZSI6ImFsb2hhNCIsImdpdmVuX25hbWUiOiJBbG9oYSIsImZhbWlseV9uYW1lIjoiVGVzdGluZyIsImVtYWlsIjoiYWxvaGE0QG1haWxpbmF0b3IuY29tIn0.i9ewOjrsRC91f08EDluC4KwOa1xSIuZNaRyKEtQP5WsQmAebFhd2EDm3J0yT5dgRbFVn5vdN4XLndfNMufpkwB9htcMxNNlZSmDKa9BdpBO0BDhQACHY23f5PJIf4GLJDiSm1rjtF3NUnFB-vGNNqH8x0ipjdtOGtIXcvb06-vSuig98racuzVvEg8WxYFWsr3m0S2xQxeC6Jr7vKbF_c25G5ymUOVySL2uT14vGKOHX3wy4yah9asy5kyuI18lBtDAo4m6xHGp4LihyVBEs7yog-tREYdjRICSGiRQyA-7c92tkfsQshJwE5h-pScetONnmN9SswfjbB_hOAnEz6Q";
    return config;
};

export const unsafeApi = axios.create({
  baseURL: "https://app-unsafe.pathwar.land/"
});

unsafeApi.interceptors.request.use(withToken);

describe('User Session API Calls', () => {
  it('should work get user session - /user/session', async () => {

    const response = await unsafeApi.get("/user/session");
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
})

describe('Organizations API Calls', () => {
  it('should work get all organizations - /organizations', () => {

  })
})
