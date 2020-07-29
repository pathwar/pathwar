/* eslint-disable react/prop-types */
import React, { useState } from "react";
import { Button, Card, Grid, Tag } from "tabler-react";
import { useQueryParam, StringParam } from "use-query-params";
import ChallengeModal from "./ChallengeModal";

const ChallengeCard = ({ challenge }) => {
  const [modalQueryId, setModalQueryId] = useQueryParam("modal", StringParam);
  const [modalOpen, setModalOpen] = useState(modalQueryId === challenge.id);

  const { flavor, id: challengeID } = challenge;
  // const isClosed = subscriptions && subscriptions[0].status === "Closed";

  const onCloseModal = function() {
    setModalOpen(false);
    setModalQueryId(undefined);
  };

  const openModal = function() {
    setModalOpen(true);
    setModalQueryId(challengeID);
  };

  return (
    <Card>
      <Card.Status color="green" side />
      <Card.Header>
        <Card.Title>{flavor.challenge.name}</Card.Title>
        <Card.Options>
          <Button color="indigo" size="sm" icon="compass" onClick={openModal}>
            View
          </Button>
        </Card.Options>
      </Card.Header>
      <Card.Body>
        <>
          <Grid.Row>
            <Grid.Col width={12}>
              <p>
                {flavor.body ||
                  `Hello Ol'salt! Try to beat the ${flavor.challenge.name} challenge.
              Heave ho!`}
              </p>
            </Grid.Col>
          </Grid.Row>
          <ChallengeModal
            open={modalOpen}
            onClose={onCloseModal}
            challengeID={challengeID}
          />
        </>
      </Card.Body>
    </Card>
  );
};

export default React.memo(ChallengeCard);
