import React, { useState } from "react";
import { Button, Form } from "tabler-react";
import { css } from "@emotion/core";

const wrapperStyle = `
  text-align: right;
`;

const ValidateCouponForm = () => {
  const [formOpen, setFormOpen] = useState(false);

  const handleFormOpen = function() {
    setFormOpen(true);
  };

  const onCouponSubmit = function(event) {
    event.preventDefault();
    alert("Comming soon...");
    setFormOpen(false);
  };

  return (
    <div
      css={css`
        ${wrapperStyle}
      `}
    >
      {!formOpen && (
        <Button color="yellow" icon="award" size="sm" onClick={handleFormOpen}>
          Validate coupon
        </Button>
      )}

      {formOpen && (
        <Form onSubmit={onCouponSubmit}>
          <Form.InputGroup
            append={
              <Button color="yellow" type="submit" size="sm">
                Validate!
              </Button>
            }
          >
            <Form.Input placeholder="Coupon code" />
          </Form.InputGroup>
        </Form>
      )}
    </div>
  );
};

export default ValidateCouponForm;
