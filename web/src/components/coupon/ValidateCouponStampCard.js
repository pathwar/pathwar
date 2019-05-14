import * as React from "react";
import { Link } from "react-router-dom";
import { StampCard } from "tabler-react";

const ValidateCouponStamp = () => {

    return (
        <StampCard
        color="yellow"
        icon="award"
        header={
          <Link to="/validate-coupon">
            <small>Validate coupon</small>
          </Link>
        }
        footer={"Get it valid!"}
      />
    )
}

export default ValidateCouponStamp;