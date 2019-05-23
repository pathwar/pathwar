var supertest = require('supertest'),
    api = supertest(process.env.API_TARGET || 'http://localhost:8000');

describe('API', () => {
    describe('/ping', () => {
        it('should return a 200 response', () => {
            return api.get('/api/ping')
                .set('Accept', 'application/json')
                .expect(200);
        });
    });
});
