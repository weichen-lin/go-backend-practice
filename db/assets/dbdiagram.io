// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

Table account {
  id uuid [pk, default: `uuid_generate_v4()`]
  owner varchar(100) [not null]
  balance decimal [not null]
  currency varchar(30) [not null]
  created_at timestampz [not null, default: `now()`]
  last_modified_at timestampz [not null, default: `now()`]

  Indexes {
    owner
  }
}

Table entries {
  id uuid [pk, default: `uuid_generate_v4()`]
  account_id uuid [ref: > account.id, not null]
  amount decimal [not null]
  created_at timestampz [not null, default: `now()`]

  Indexes {
    account_id
  }
}

Table transfers {
  id uuid [pk, default: `uuid_generate_v4()`]
  from_account_id uuid [ref: > account.id, not null]
  to_account_id uuid [ref: > account.id, not null]
  amount decimal [not null]
  created_at timestampz [not null, default: `now()`]

  Indexes {
    from_account_id
    to_account_id
    (from_account_id, to_account_id)
  }
}



