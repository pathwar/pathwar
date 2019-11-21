import * as React from "react";
import { Link } from "gatsby";
import { StampCard } from "tabler-react";

const CreateTeamStampCard = () => {

    return (
        <StampCard
        color="blue"
        icon="users"
        header={
          <Link to="/validate-coupon">
            <small>Create new team</small>
          </Link>
        }
        footer={"Ahoy! Create a new team"}
      />
    )
}

export default CreateTeamStampCard;
