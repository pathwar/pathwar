import React, {useEffect, useState} from "react";
import {isEmpty} from "ramda";
import {Button, Form} from "tabler-react";
import {FormattedMessage} from "react-intl";

const styles = {
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
}

const InviteOrganizationButton = () => {

  const [isFormOpen, setFormOpen] = useState(false);
  const [name, setName] = useState("");
  const [error, setError] = useState(false);

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
        <form>
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
