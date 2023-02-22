import React from "react";
import {Avatar} from "tabler-react";
import {Link} from "gatsby";

export const UserOrganizationBadges = ({organizations}) => {
  return !organizations ? (
    <></>
    ) : (
      <Avatar.List>
        {organizations.map(organization => (
          <Link key={organization.id} to={"/organization/" + organization.id}>
            <Avatar size="md" imageURL={`${organization.gravatar_url}?d=identicon`}/>
          </Link>
     ))}
      </Avatar.List>
  );
}

export default UserOrganizationBadges;
