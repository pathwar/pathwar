/* eslint-disable react/prop-types */
import React from "react";
import { Modal } from "react-responsive-modal";
import ChallengeDetailsPage from "../../pages/ChallengeDetailsPage";

const ChallengeModal = ({ open, onClose, challengeID }, ...rest) => {
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
