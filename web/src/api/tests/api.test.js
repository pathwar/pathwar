import axios from "axios";
require('dotenv').config({ path: "../../.env" });

//Axios Config
const withToken = function(config) {
    config.headers.Authorization = process.env.KEYCLOAK_TOKEN;
    return config;
};

export const unsafeApi = axios.create({
  baseURL: process.env.API_URL_UNSAFE
});

unsafeApi.interceptors.request.use(withToken);

//Test session variables
let active_season_id = undefined
let active_team_id = undefined
let season_challenge_id = undefined

let challenge_subscription_id = undefined

//Helpers
const performUserSessionCalls = async () => {

  // ensure we have an account for this token before deleting it in the next step
  try {
    const userSessionResponse = await unsafeApi.get("/user/session");
    const { user } = userSessionResponse.data
  } catch (error) {
    throw error;
  }

  // trash any existing account first
  try {
    const response = await unsafeApi.post(`/user/delete-account`, {"reason": "integration test"})
  } catch (error) {
    throw error;
  }

  // Set a real season id and team id for tests
  try {
    const userSessionResponse = await unsafeApi.get("/user/session");
    const { user } = userSessionResponse.data

    active_season_id = user.active_season_id
    active_team_id = user.active_team_member.team_id

    console.log("Season ID >>", active_season_id)
    console.log("Active team ID >>", active_team_id)
  } catch (error) {
    throw error;
  }

  // Set a real challenge_id from the season for tests
  try {
    const seasonChallengesResponse = await unsafeApi.get(`/season-challenges?season_id=${active_season_id}`);
    const firstItem = seasonChallengesResponse.data.items[0]
    season_challenge_id = firstItem.id
    console.log("Challenge ID >>", season_challenge_id)
  } catch (error) {
    throw error;
  }

}

beforeAll((done) => {
  jest.setTimeout(50000);
  performUserSessionCalls();
  return done();
});

describe('API Calls', () => {
  it('should work GET user session - /user/session', async () => {
    const response = await unsafeApi.get("/user/session");
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
  it('should work GET all organizations - /organizations', async () => {
    const response = await unsafeApi.get("/organizations");
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
  it('should work POST preferences - /user/preferences', async () => {
    const preferencesPost = await unsafeApi.post(`/user/preferences`, {"active_season_id": active_season_id});
    expect(preferencesPost.status).toEqual(200);
  })
  it('should work GET all teams on a season - /teams?season_id=the_id', async () => {
    const getAllTeamsResponse = await unsafeApi.get(`/teams?season_id=${active_season_id}`);
    expect(getAllTeamsResponse.status).toEqual(200);
  })
  it('should work GET all challenges on a season - /season-challenges', async () => {
    const response = await unsafeApi.get(`/season-challenges?season_id=${active_season_id}`);
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
  it('should work GET season challenge details - /season-challenge?season_challenge_id=the_id', async () => {
    const response = await unsafeApi.get(`/season-challenge?season_challenge_id=${season_challenge_id}`);
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
  it('should work GET team details - /team?team_id=the_id', async () => {
    const response = await unsafeApi.get(`/team?team_id=${active_team_id}`);
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
  it('should work POST season challenge BUY - /season-challenge/buy', async() => {
    const response = await unsafeApi.post(`/season-challenge/buy`, {"season_challenge_id": season_challenge_id, "team_id": active_team_id});
    const { challenge_subscription } = response.data;
    challenge_subscription_id = challenge_subscription.id;
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
  it('should work POST season challenge VALIDATE - /season-challenge/validate', async() => {
    const response = await unsafeApi.post(`/season-challenge/validate`, {"challenge_subscription_id": challenge_subscription_id,  "passphrase": "lorem ipsum", "comment": "dolor sit amet"});
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
  it('should work POST season challenge CLOSE - /season-challenge/close', async() => {
    const response = await unsafeApi.post(`/season-challenge/close`, {"challenge_subscription_id": challenge_subscription_id});
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
})
