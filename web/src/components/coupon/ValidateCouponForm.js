import React, { useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Button, Form } from "tabler-react";
import { css } from "@emotion/core";
import { FormattedMessage } from "react-intl";
import { fetchCouponValidation } from "../../actions/userSession";

const wrapperStyle = `
  text-align: right;
  margin-right: 1rem;
  display: inline-block;
`;

const ValidateCouponForm = () => {
  const dispatch = useDispatch();
  const [formOpen, setFormOpen] = useState(false);
  const [code, setCode] = useState("");
  const activeTeam = useSelector(state => state.seasons.activeTeam);

  const handleFormOpen = function() {
    setFormOpen(true);
  };

  const handleChange = event => {
    setCode(event.target.value);
  };

  const onCouponSubmit = function(event) {
    event.preventDefault();
    dispatch(fetchCouponValidation(code, activeTeam.id));
    setFormOpen(false);
    setCode("");
  };

  return (
    <div
      css={css`
        ${wrapperStyle}
      `}
    >
      {!formOpen && (
        <Button color="yellow" icon="award" size="sm" onClick={handleFormOpen}>
          <FormattedMessage id="coupon" />
        </Button>
      )}

      {formOpen && (
        <Form onSubmit={onCouponSubmit}>
          <Form.InputGroup
            append={
              <Button color="yellow" type="submit" size="sm">
                <FormattedMessage id="coupon.validate" />
              </Button>
            }
          >
            <Form.Input
              onChange={handleChange}
              placeholder="Coupon code"
              autoFocus={true}
            />
          </Form.InputGroup>
        </Form>
      )}
    </div>
  );
};

export default ValidateCouponForm;
