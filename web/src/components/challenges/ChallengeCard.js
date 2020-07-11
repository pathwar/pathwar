/* eslint-disable react/prop-types */
import React, { useState } from "react";
import { Button, Card, Grid, Tag } from "tabler-react";
import { Modal } from "react-responsive-modal";
import ChallengeDetailsPage from "../../pages/ChallengeDetailsPage";

const ChallengeCard = ({ challenge }) => {
  const [modalOpen, setModalOpen] = useState(false);
  const { flavor, id: challengeID } = challenge;
  // const isClosed = subscriptions && subscriptions[0].status === "Closed";

  const onCloseModal = function() {
    setModalOpen(false);
  };

  const openModal = function() {
    setModalOpen(true);
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
            <Grid.Col auto>
              <p>Author: {flavor.challenge.author}</p>
            </Grid.Col>
          </Grid.Row>
          <Grid.Row>
            <Grid.Col auto>
              <h3>changelog:</h3>
              <p>{flavor.changelog}</p>
            </Grid.Col>
          </Grid.Row>
          <Grid.Row>
            <Grid.Col auto>
              <Tag.List>
                <Tag color="dark" addOn={flavor.version} addOnColor="warning">
                  version
                </Tag>
                <Tag addOn={flavor.is_latest.toString()} addOnColor="success">
                  is_latest
                </Tag>
              </Tag.List>
            </Grid.Col>
          </Grid.Row>
          <Modal open={modalOpen} onClose={onCloseModal}>
            <ChallengeDetailsPage challengeID={challengeID} />
          </Modal>
        </>
      </Card.Body>
    </Card>
  );
};

export default React.memo(ChallengeCard);
