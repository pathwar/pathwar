import * as React from "react";
import { Link } from "react-router-dom";
import { StampCard } from "tabler-react";

const TeamStatsStampCard = () => {

    return (
        <StampCard
        color="red"
        icon="activity"
        header={
          <Link to="/statistics">
            <small>See Teams Stats</small>
          </Link>
        }
        footer={"List teams statistics"}
      />
    )
}

export default TeamStatsStampCard;