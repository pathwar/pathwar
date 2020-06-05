import React from "react";
import { Card, Grid } from "tabler-react";

const ValidationsList = ({ validations }) => {
  return (
    <Grid.Row cards={true}>
      {validations.map(validation => {
        const status = validation.status;
        const statusColor =
          status === "NeedReview"
            ? "orange"
            : status === "Rejected"
            ? "red"
            : "green";
        return (
          <Grid.Col lg={4} md={4} sm={6} xs={6} key={validation.id}>
            <Card
              title={validation.passphrases}
              statusColor={statusColor}
              isCollapsible
            >
              <Card.Body>
                <p>
                  <strong>Status:</strong> {validation.status}
                </p>
                <p>
                  <strong>Key:</strong> {validation.passphrase_key}
                </p>
                <p>
                  <strong>Comment:</strong> {validation.author_comment}
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
