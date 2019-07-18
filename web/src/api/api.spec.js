var supertest = require('supertest'),
    api = supertest(process.env.API_TARGET || 'http://localhost:8000');

describe('API', () => {
    describe('/ping', () => {
        it('should return a 200 response', () => {
            return api.get('/api/ping')
                .set('Accept', 'application/json')
                .expect(200);
        });
        it('should return a 501 response', () => {
            return api.post('/api/ping')
                .set('Accept', 'application/json')
                .expect(501);
        });
    });
    var token = '';
    describe('login flow', () => {
        it('should return a token', () => {
            return api.post('/api/authenticate')
                .send({'username': 'integration'})
                .expect(200)
                .expect(response => {
                    if (!('token' in response.body)) {
                        throw new Error('missing token');
                    }
                })
                .then(response => {
                    token = response.body.token;
                })
        });
        it('should return a 401 response without token', () => {
            return api.get('/api/user-session')
                .expect(401);
        });
        it('should return a 200 response with a valid token', () => {
            return api.get('/api/user-session')
                .set('Authorization', token)
                .expect(200, {
                    'metadata': {},
                    'username': 'integration',
                });
        });
    });
});
