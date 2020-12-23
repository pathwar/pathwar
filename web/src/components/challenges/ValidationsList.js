import React from "react";
import { Card, Grid } from "tabler-react";
import { css } from "@emotion/core";
import { FormattedMessage } from "react-intl";

const validationCard = css`
  text-align: left;
  margin-top: 0.5rem;
`;

const title = css`
  margin-top: 1rem;
`;

const ValidationsList = ({ validations }) => {
  return (
    <Grid.Row cards={true}>
      <Grid.Col width={12}>
        <h3 css={title}>
          <FormattedMessage id="ValidationList.title" />:
        </h3>
      </Grid.Col>
      {validations.map(validation => {
        const status = validation.status;
        const statusColor =
          status === "NeedReview"
            ? "orange"
            : status === "Rejected"
            ? "red"
            : "green";
        return (
          <Grid.Col width={12} sm={6} md={4} key={validation.id}>
            <Card
              css={validationCard}
              title={validation.passphrases}
              statusColor={statusColor}
              isCollapsible
            >
              <Card.Body>
                <p>
                  <strong>
                    <FormattedMessage id="ValidationList.status" />:
                  </strong>{" "}
                  {validation.status}
                </p>
                <p>
                  <strong>
                    <FormattedMessage id="ValidationList.key" />:
                  </strong>{" "}
                  {validation.passphrase_key}
                </p>
                <p>
                  <strong>
                    <FormattedMessage id="ValidationList.comment" />:
                  </strong>{" "}
                  {validation.author_comment}
                </p>
              </Card.Body>
            </Card>
          </Grid.Col>
        );
      })}
    </Grid.Row>
  );
};

export default ValidationsList;
