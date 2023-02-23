import React, {useEffect, useState} from "react";
import {isEmpty} from "ramda";
import {Button, Form} from "tabler-react";
import {FormattedMessage} from "react-intl";
import {inviteUserToOrganization as InviteUserToOrganizationAction} from "../../actions/organizations";
import {useDispatch} from "react-redux";

const styles = {
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
}

const InviteOrganizationButton = ({organizationID, organizationName}) => {
  const dispatch = useDispatch();
  const [isFormOpen, setFormOpen] = useState(false);
  const [name, setName] = useState("");
  const [error, setError] = useState(false);
  const inviteUserToOrganization = (organizationID, name, organizationName) =>
    dispatch(InviteUserToOrganizationAction(organizationID, name, organizationName));

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

  const submitOrganizationInvite = async event => {
    event.preventDefault();
    if (isEmpty(name)) {
      setError(true);
    } else {
      await inviteUserToOrganization(organizationID, name, organizationName);
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
        {"Invite new member"}
      </Button>
      {isFormOpen && (
        <form onSubmit={submitOrganizationInvite}>
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
              <Button type="submit" color="primary" className="ml-auto">
                <FormattedMessage id="CreateTeamButton.send" />
              </Button>
            </Form.Group>
          </Form.FieldSet>
        </form>
      )}
    </div>
  );
}

export default InviteOrganizationButton;
