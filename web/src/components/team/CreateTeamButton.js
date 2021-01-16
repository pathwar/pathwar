import React, { useState, useEffect } from "react";
import { isEmpty } from "ramda";
import { Form, Button } from "tabler-react";
import { FormattedMessage } from "react-intl";

const CreateTeamButton = ({ activeSeason, activeTeamInSeason, createTeam }) => {
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

  const submitTeamCreate = async event => {
    event.preventDefault();
    if (isEmpty(name)) {
      setError(true);
      return;
    } else {
      await createTeam(activeSeason.id, name);
      setFormOpen(false);
    }
  };

  return (
    <>
      <Button
        color="success"
        onClick={handleFormOpen}
        icon={activeTeamInSeason ? "anchor" : "users"}
        disabled={activeTeamInSeason}
        size="sm"
      >
        {activeTeamInSeason ? "Team on season" : "Create team"}
      </Button>
      {isFormOpen && (
        <form onSubmit={submitTeamCreate}>
          <Form.FieldSet>
            <Form.Group isRequired label="Name">
              <Form.Input
                name="name"
                onChange={handleChange}
                invalid={error}
                cross={error}
                feedback={error && "Please, insert a name"}
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
    </>
  );
};

export default CreateTeamButton;
