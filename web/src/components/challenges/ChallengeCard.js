/* eslint-disable react/prop-types */
import React, { useState } from "react";
import { useQueryParam, StringParam } from "use-query-params";
import { FormattedMessage } from "react-intl";
import { css } from "@emotion/core";
import ChallengeModal from "./ChallengeModal";
import ChallengeBuyButton from "./ChallengeBuyButton";
import iconPwn from "../../images/icon-pwn-small.svg";
import generateLightColorHex from "../../utils/generateLightColorHex";

const mainContainer = css`
  margin-bottom: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
`;

const cardWrapper = isClosed => css`
  background-color: #fff;
  display: flex;
  flex-direction: row;
  align-items: center;
  box-shadow: 0px 5px 20px 0px rgba(7, 42, 68, 0.1);
  margin-bottom: 0.5rem;
  padding: 1rem 1rem;
  border-radius: 8px;
  min-height: 374px;
  width: 100%;
  ${isClosed && `pointer-events: none`}
`;

const cardActionsWrapper = css`
  min-width: 206px;
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
      font-size: 1.125rem;
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

      .boldItem {
        margin: 0;
        font-weight: bold;
      }
    }
  }
`;

const tagStyle = bgColor => css`
  display: inline-block;
  text-align: center;
  background-color: ${bgColor};
  border-radius: 8px;
  padding: 0.5rem;
  margin-right: 0.5rem;
  font-weight: 300;
  font-size: 0.75rem;
`;

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

  return (
    <>
      <div css={mainContainer} onClick={openModal}>
        <div css={() => cardWrapper(isClosed)}>
          <div css={cardActionsWrapper}>
            <img
              src="https://d33wubrfki0l68.cloudfront.net/1c254da613f195cbfc2a85e94c1f792b306abea4/09aac/files/islands--pathwar-island-desert.svg"
              alt="island"
              width={206}
              height={206}
            />
            {!isClosed && !subscriptions && (
              <div>
                <ChallengeBuyButton challenge={challenge} />
              </div>
            )}
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
                <div css={theme => tagStyle(theme.colors.gray)}>
                  {flavor.category}
                </div>
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
                {!subscriptions && (
                  <p className="boldItem">
                    <FormattedMessage id="ChallengeCard.price" />:{" "}
                    {flavor.purchase_price ? `$${flavor.purchase_price}` : "$0"}{" "}
                    <img src={iconPwn} />
                  </p>
                )}
                <p className="boldItem">
                  <FormattedMessage id="ChallengeCard.reward" />:{" "}
                  {`$${flavor.validation_reward}`} <img src={iconPwn} />
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
      <ChallengeModal
        open={modalOpen}
        onClose={onCloseModal}
        challengeID={challengeID}
      />
    </>
  );
};

export default React.memo(ChallengeCard);
