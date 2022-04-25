import * as R from "ramda";
import { getAllSeasonTeams } from "../api/seasons";

const getTeamRank = async (teamId, seasonId) => {
  const scoreCashSort = R.sortWith([
    R.descend(R.prop("score")),
    R.descend(R.prop("cash")),
  ]);

  try {
    const response = await getAllSeasonTeams(seasonId);
    const teams = response.data.items;

    const parsedTeams = teams.map(item => ({
      ...item,
      score: item.score ? parseInt(item.score, 10) : undefined,
      cash: item.cash ? parseInt(item.cash, 10) : undefined,
    }));

    const sortedTeamsByScoreAndCash = scoreCashSort(parsedTeams);
    const rank = sortedTeamsByScoreAndCash.findIndex(
      item => item.id === teamId
    );

    return rank <= 0 ? 0 : rank;
  } catch (e) {
    return e;
  }
};

export default getTeamRank;
