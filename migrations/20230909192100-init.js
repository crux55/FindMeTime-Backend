'use strict';

var dbm;
var type;
var seed;

exports.setup = function (options, seedLink) {
  dbm = options.dbmigrate;
  type = dbm.dataType;
  seed = seedLink;
};

exports.up = function (db) {
  return db.createTable('users', {
    columns: {
      id: { type: 'VARCHAR(50)', primaryKey: true },
      username: { type: 'VARCHAR(20)', notNull: true },
    },
    ifNotExists: true
  })
  .then(() => {
    return db.createTable('tasks', {
      columns: {
        id: { type: 'VARCHAR(50)', primaryKey: true },
        title: { type: 'VARCHAR(50)', notNull: true },
        description: { type: 'VARCHAR(50)', notNull: true },
        duration: { type: 'int', notNull: true },
        created_on: { type: 'timestamp', notNull: true },
        tags_only: { type: 'VARCHAR(5000)' },
        tags_not: { type: 'VARCHAR(5000)', arrayType: 'string' },
      },
      ifNotExists: true
    });
  })
  .then(() => {
    return db.createTable('tags', {
      columns: {
        id: { type: 'VARCHAR(50)', primaryKey: true },
        tag_name: { type: 'VARCHAR(20)', notNull: true },
        description: { type: 'VARCHAR(20)', notNull: true },
        time_slots: { type: 'VARCHAR(5000)' },
      },
      ifNotExists: true
    });
  })
  .then(() => {
    return db.createTable('time_slots', {
      columns: {
        id: { type: 'VARCHAR(50)', primaryKey: true },
        start_day_index: { type: 'int' },
        start_time: { type: 'int' },
        end_day_index: { type: 'int' },
        end_time: { type: 'int' },
      },
      ifNotExists: true
    });
  })
  .then(() => {
    db.runSql('CREATE ROLE tasker')
  })
  .then(() => {
    // Grant permissions to 'tasker' user or role.
    return db.runSql('GRANT ALL ON TABLE users TO tasker');
  })
  .then(() => {
    return db.runSql('GRANT ALL ON TABLE tasks TO tasker');
  })
  .then(() => {
    return db.runSql('GRANT ALL ON TABLE tags TO tasker');
  })
  .then(() => {
    return db.runSql('GRANT ALL ON TABLE time_slots TO tasker');
  });
};

exports.down = function (db) {
  return db.dropTable('time_slots')
    .then(() => db.dropTable('tags'))
    .then(() => db.dropTable('tasks'))
    .then(() => db.dropTable('users'));
};

exports._meta = {
  version: 1,
};