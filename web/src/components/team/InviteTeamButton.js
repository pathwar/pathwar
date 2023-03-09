import React, {useEffect, useState} from "react";
import {isEmpty} from "ramda";
import {Button, Form} from "tabler-react";
import {FormattedMessage} from "react-intl";
import {inviteUserToTeam as inviteUserToTeamAction} from "../../actions/seasons";
import {useDispatch} from "react-redux";

const styles = {
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  marginBottom: "0.5rem",
}

const InviteTeamButton = ({teamID, seasonName, organizationName}) => {
  const dispatch = useDispatch();
  const [isFormOpen, setFormOpen] = useState(false);
  const [name, setName] = useState("");
  const [error, setError] = useState(false);
  const inviteUserToTeam = (teamID, name, organizationName, seasonName) =>
    dispatch(inviteUserToTeamAction(teamID, name, organizationName, seasonName));

  useEffect(() => {
    if (!isEmpty(name)) {
      setError(false);
    }
  }, [name]);

  const handleChange = event => setName(event.target.value);
  const handleFormOpen = event => {
    event.preventDefault();
    setFormOpen(!isFormOpen);
  };

  const submitTeamInvite = async event => {
    event.preventDefault();
    if (isEmpty(name)) {
      setError(true);
    } else {
      await inviteUserToTeam(teamID, name, organizationName, seasonName);
      setFormOpen(false);
    }
  }

  return (
    <div style={styles}>
      <Button
        color="success"
        onClick={handleFormOpen}
        icon={"users"}
        size="sm"
        css={styles}
      >
        {"Invite new member in team"}
      </Button>
      {isFormOpen && (
        <form onSubmit={submitTeamInvite}>
          <Form.FieldSet css={styles}>
            <Form.Group isRequired label="Username" css={styles}>
              <Form.Input
                name="name"
                onChange={handleChange}
                invalid={error}
                cross={error}
                feedback={error && "Please, insert a username"}
              />
            </Form.Group>
            <Form.Group>
              <Button type="submit" color="success" className="ml-auto">
                <FormattedMessage id="CreateTeamButton.send" />
              </Button>
            </Form.Group>
          </Form.FieldSet>
        </form>
      )}
    </div>
  );
}

export default InviteTeamButton;
