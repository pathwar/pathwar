/* eslint-disable react/prop-types */
import React from "react";
import { Modal } from "react-responsive-modal";
import ChallengeDetailsPage from "../../pages/ChallengeDetailsPage";

const ChallengeModal = (
  { open, onClose: onCloseProps, challengeID },
  ...rest
) => {
  return (
    <Modal
      open={open}
      onClose={onCloseProps}
      center={true}
      animationDuration={300}
      {...rest}
    >
      <ChallengeDetailsPage challengeID={challengeID} />
    </Modal>
  );
};

export default React.memo(ChallengeModal);
