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
let team_id = undefined
let season_challenge_id = undefined

//Helpers
const performUserSessionCalls = async () => {

  //Set a real season id for tests
  try {
    const userSessionResponse = await unsafeApi.get("/user/session");
    const { user } = userSessionResponse.data
    active_season_id = user.active_season_id
  } catch (error) {
    throw error;
  }

  //Set a real challenge_id from the season for tests
  try {
    const seasonChallengesResponse = await unsafeApi.get(`/season-challenges?season_id=${active_season_id}`);
    const firstItem = seasonChallengesResponse.data.items[0]
    season_challenge_id = firstItem.id
  } catch (error) {
    throw error;
  }

  //Set a real team_id from the teams on season
  try {
    const getAllTeamsResponse = await unsafeApi.get(`/teams?season_id=${active_season_id}`);
    const { items } = getAllTeamsResponse.data;
    const team = items.find((item) => item.is_default)
    team_id = team.id;
  } catch (error) {
    throw error;
  }

}

beforeAll(() => {
  return performUserSessionCalls();
})

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
  it('should work POST preferences - /preferences', async () => {
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
    const response = await unsafeApi.get(`/team?team_id=${team_id}`);
    expect(response.status).toEqual(200);
    expect(response.data).toBeDefined();
  })
})
