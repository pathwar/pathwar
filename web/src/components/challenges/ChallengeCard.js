/* eslint-disable react/prop-types */
import React, { useState } from "react";
import { Card, Grid, Tag } from "tabler-react";
import { useQueryParam, StringParam } from "use-query-params";
import { FormattedMessage } from "react-intl";
import { css } from "@emotion/core";
import ChallengeModal from "./ChallengeModal";
import Button from "../Button";
import ChallengeBuyButton from "./ChallengeBuyButton";

const cardTag = css`
  margin-right: 0.5rem;
`;

const notPurchased = css`
  box-shadow: none;
  opacity: 0.8;
`;

const closedStyle = css`
  opacity: 0.4;
  pointer-events: none;
  box-shadow: none;
`;

const cardWrapper = css`
  background-color: #fff;
  display: flex;
  flex-direction: row;
  box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
  margin-bottom: 1.5rem;
  padding: 3rem 1rem;
  border-radius: 8px;
`;

const cardActionsWrapper = theme => css`
  display: flex;
  flex-direction: column;
  flex: 1;

  .pricingWrapper {
    display: flex;
  }

  .pill {
    border-radius: 5px;
    padding: 0.5rem;
    text-align: center;
    font-weight: bold;
    background-color: ${theme.colors.success};
    box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
    text-transform: uppercase;

    width: 100%;
  }
`;

const ChallengeCard = ({ challenge }) => {
  const [modalQueryId, setModalQueryId] = useQueryParam("modal", StringParam);
  const [modalOpen, setModalOpen] = useState(modalQueryId === challenge.id);

  const { flavor, subscriptions, id: challengeID } = challenge;
  const purchased = subscriptions;
  const isClosed = subscriptions && subscriptions[0].status === "Closed";

  const onCloseModal = () => {
    setModalOpen(false);
    setModalQueryId(undefined);
  };

  const openModal = () => {
    setModalOpen(true);
    setModalQueryId(challengeID);
  };

  const cardColor = isClosed ? "red" : purchased ? "blue" : "gray";

  return (
    <>
      <div css={cardWrapper}>
        <div css={theme => cardActionsWrapper(theme)}>
          <img
            src="https://d33wubrfki0l68.cloudfront.net/1c254da613f195cbfc2a85e94c1f792b306abea4/09aac/files/islands--pathwar-island-desert.svg"
            alt="island"
          />
          <div className="pricingWrapper">
            <div className="pill">
              <FormattedMessage id="ChallengeCard.price" />
            </div>
            <div className="pill">
              {flavor.purchase_price ? `$${flavor.purchase_price}` : "$0"}
            </div>
          </div>
          <div>
            <ChallengeBuyButton challenge={challenge} />
          </div>
        </div>
        <div>DESCRIPTION</div>
      </div>
      <ChallengeModal
        open={modalOpen}
        onClose={onCloseModal}
        challengeID={challengeID}
      />
      {/* <Card css={[!purchased && notPurchased, isClosed && closedStyle]}>
        <Card.Status color={cardColor} side />
        <Card.Header>
          <Card.Title>{flavor.challenge.name}</Card.Title>
          <Card.Options>
            <Button
              color={isClosed ? "red" : "indigo"}
              size="sm"
              icon={isClosed ? "x-circle" : "compass"}
              onClick={openModal}
            >
              {isClosed ? (
                <FormattedMessage id="ChallengeCard.closed" />
              ) : (
                <FormattedMessage id="ChallengeCard.view" />
              )}
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
            <Grid.Row>
              <Grid.Col width={12}>
                <Tag
                  css={cardTag}
                  color="lime"
                  addOn={
                    flavor.purchase_price ? `$${flavor.purchase_price}` : "$0"
                  }
                  addOnColor="green"
                >
                  <FormattedMessage id="ChallengeCard.price" />
                </Tag>
                <Tag
                  css={cardTag}
                  color="yellow"
                  addOn={flavor.validation_reward}
                  addOnColor="yellow"
                >
                  <FormattedMessage id="ChallengeCard.reward" />
                </Tag>
              </Grid.Col>
            </Grid.Row>
            <ChallengeModal
              open={modalOpen}
              onClose={onCloseModal}
              challengeID={challengeID}
            />
          </>
        </Card.Body>
      </Card> */}
    </>
  );
};

export default React.memo(ChallengeCard);
