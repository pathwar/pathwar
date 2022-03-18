/* eslint-disable react/prop-types */
import React, { useState } from "react";
import { useQueryParam, StringParam } from "use-query-params";
import { FormattedMessage } from "react-intl";
import { css } from "@emotion/core";
import ChallengeModal from "./ChallengeModal";
import Button from "../Button";
import ChallengeBuyButton from "./ChallengeBuyButton";
import mapIcon from "../../images/map-l-icon.svg";
import iconPwn from "../../images/icon-pwn-small.svg";
import generateLightColorHex from "../../utils/generateLightColorHex";

const mainContainer = css`
  margin-bottom: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;

  .infoRow {
    position: absolute;
    top: 0.5rem;
    right: 1.5rem;
  }
`;

const cardWrapper = isClosed => css`
  background-color: #fff;
  display: flex;
  flex-direction: row;
  box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
  margin-bottom: 0.5rem;
  padding: 3rem 1rem;
  border-radius: 8px;
  max-height: 374px;
  width: 100%;
  ${isClosed && `pointer-events: none`}
`;

const cardActionsWrapper = css`
  display: flex;
  flex-direction: column;
  min-width: 206px;

  .pricingWrapper {
    display: flex;
  }
`;

const cardDescriptionWrapper = theme => css`
  padding: 0 1rem;
  flex: 1;

  .descriptionHeader {
    padding: 0;
    border: none;
    margin-bottom: 1rem;

    .title,
    .subtitle {
      font-size: 1.313rem;
      margin: 0;
    }

    .title {
      font-weight: bolder;
      border-bottom: 1px solid;
      padding-bottom: 1rem;
    }

    .subtitle {
      font-weight: bold;
      padding: 1rem 0 0 0;
    }
  }

  .descriptionBody {
    .tagsWrapper {
      margin-bottom: 0.5rem;
    }
    .statsWrapper {
      font-size: 0.75rem;
      padding: 1rem;
      border-radius: 8px;
      background-color: ${theme.colors.gray};
      box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);

      .heading {
        text-align: center;
        font-weight: bold;
      }

      .item {
        margin-bottom: 0.2rem;
        span {
          font-weight: bold;
        }
      }

      .rewardItem {
        margin: 0;
        font-weight: bold;
      }
    }
  }
`;

const pillStyle = color => css`
  border-radius: 5px;
  padding: 0.5rem;
  text-align: center;
  font-weight: 900;
  background-color: ${color};
  box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
  text-transform: uppercase;

  width: 100%;

  &:nth-of-type(2) {
    margin-left: 0.2rem;
  }
`;

const tagStyle = bgColor => css`
  display: inline-block;
  text-align: center;
  background-color: ${bgColor};
  border-radius: 8px;
  padding: 0.5rem;
  font-weight: 300;
  font-size: 0.75rem;
`;

const islandsUrls = [
  "https://d33wubrfki0l68.cloudfront.net/1c254da613f195cbfc2a85e94c1f792b306abea4/09aac/files/islands--pathwar-island-desert.svg",
  "https://d33wubrfki0l68.cloudfront.net/055daf3aaf80e0bacca1db87ad3ffa001d294c69/82028/files/islands--pathwar-island-grassland.svg",
  "https://d33wubrfki0l68.cloudfront.net/1eaa050952258bfe93771ef87c9edb9c990b9ba3/31cd1/files/islands--pathwar-island-jungle.svg",
  "https://d33wubrfki0l68.cloudfront.net/0012a2d6917afcc1ace2e0363eef0f4697aa04b7/38329/files/islands--pathwar-island-mountain.svg",
  "https://d33wubrfki0l68.cloudfront.net/81a67e12a8a0a40d1148465e8a777a25d1cb07bf/25190/files/islands--pathwar-island-north.svg",
  "https://d33wubrfki0l68.cloudfront.net/7bdd23eb7ad738d1a00122d20e50c4a0df46dd2d/21303/files/islands--pathwar-island-volcano.svg",
];

const ChallengeCard = ({ challenge }) => {
  const [modalQueryId, setModalQueryId] = useQueryParam("modal", StringParam);
  const [modalOpen, setModalOpen] = useState(modalQueryId === challenge.id);

  const { flavor, subscriptions, id: challengeID } = challenge;
  const isClosed = subscriptions && subscriptions[0].status === "Closed";

  const onCloseModal = () => {
    setModalOpen(false);
    setModalQueryId(undefined);
  };

  const openModal = () => {
    setModalOpen(true);
    setModalQueryId(challengeID);
  };

  const islandImg = islandsUrls[(Math.random() * islandsUrls.length) | 0];

  return (
    <div css={mainContainer}>
      <div className="infoRow">
        <div
          css={theme => [
            pillStyle(
              isClosed ? theme.colors.secondary : theme.colors.darkGreen
            ),
            `color: ${theme.colors.light}; font-size: 0.75rem`,
          ]}
        >
          {isClosed ? (
            <FormattedMessage id="ChallengeCard.closed" />
          ) : (
            <FormattedMessage id="ChallengeCard.open" />
          )}
        </div>
      </div>
      <div css={() => cardWrapper(isClosed)}>
        <div css={cardActionsWrapper}>
          <img src={islandImg} alt="island" width={206} height={206} />
          <div className="pricingWrapper">
            <div css={theme => pillStyle(theme.colors.success)}>
              <FormattedMessage id="ChallengeCard.price" />
            </div>
            <div css={theme => pillStyle(theme.colors.success)}>
              {flavor.purchase_price ? `$${flavor.purchase_price}` : "$0"}
            </div>
          </div>
          <div>
            <ChallengeBuyButton challenge={challenge} />
          </div>
        </div>
        <div css={theme => cardDescriptionWrapper(theme)}>
          <div className="descriptionHeader">
            <h2 className="title">{flavor.challenge.name}</h2>
            <h2 className="subtitle">
              {flavor.body ||
                `Hello Ol'salt! Try to beat the ${flavor.challenge.name} challenge.
              Heave ho!`}
            </h2>
          </div>
          <div className="descriptionBody">
            <div className="tagsWrapper">
              {flavor.tags.map(tag => (
                <div css={tagStyle(generateLightColorHex())} key={tag}>
                  {tag}
                </div>
              ))}
            </div>
            <div className="statsWrapper">
              <p className="heading">
                <FormattedMessage id="ChallengeCard.statsHeading" />:
              </p>
              <p className="item">
                <FormattedMessage id="ChallengeCard.category" />:{" "}
                <span>{flavor.category}</span>
              </p>
              <p className="rewardItem">
                {flavor.validation_reward} <img src={iconPwn} />{" "}
                <FormattedMessage id="ChallengeCard.reward" />
              </p>
            </div>
          </div>
        </div>
      </div>
      <Button color="secondary" onClick={openModal}>
        <img width={35} height={35} src={mapIcon} />{" "}
        <FormattedMessage id="ChallengeCard.view" />!
      </Button>
      <ChallengeModal
        open={modalOpen}
        onClose={onCloseModal}
        challengeID={challengeID}
      />
    </div>
  );
};

export default React.memo(ChallengeCard);
