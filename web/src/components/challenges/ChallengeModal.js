/* eslint-disable react/prop-types */
import React from "react";
import { Modal } from "react-responsive-modal";
import ChallengeDetailsPage from "../../pages/ChallengeDetailsPage";

const browser = typeof window !== "undefined" && window;

const ChallengeModal = (
  { open, onClose: onCloseProps, challengeID },
  ...rest
) => {
  if (open) {
    browser && window.history.pushState(null, null, `?modal=${challengeID}`);
  }

  const onClose = () => {
    browser && window.history.pushState(null, null, "/app/challenges");
    onCloseProps && onCloseProps();
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      center={true}
      animationDuration={300}
      {...rest}
    >
      <ChallengeDetailsPage challengeID={challengeID} />
    </Modal>
  );
};

export default React.memo(ChallengeModal);
