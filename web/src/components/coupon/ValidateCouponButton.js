import * as React from "react";
import { Button } from "tabler-react";

const ValidateCouponButton = () => {
  return (
    <Button
      color="yellow"
      icon="award"
      size="sm"
      onClick={() => alert("Check this feature!")}
    >
      Validate coupon
    </Button>
  );
};

export default ValidateCouponButton;
