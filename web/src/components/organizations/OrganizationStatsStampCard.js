import * as React from "react";
import { Link } from "gatsby";
import { StampCard } from "tabler-react";

const OrganizationStatsStampCard = () => {

    return (
        <StampCard
        color="red"
        icon="activity"
        header={
          <Link to="/app/statistics">
            <small>See Teams Stats</small>
          </Link>
        }
        footer={"List teams statistics"}
      />
    )
}

export default OrganizationStatsStampCard;
