/*
проект AirBnB
онлайн-площадка для размещения, поиска и аренды жилья.
Содержит данные о субъектах(арендаторы и арендодатели), объектах аренды и контрактах.
*/



-- Субъекты

CREATE TABLE public.subjects (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	firstname text NOT NULL,
	middlename text NULL,
	lastname text NULL,
	email text NULL,
	phone text NULL,
	userpic text NULL,
	CONSTRAINT subjects_pk PRIMARY KEY (id)
);


-- Счета

CREATE TABLE public.accounts (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	account_number varchar NOT NULL,
	id_subject bigint NOT NULL,
	CONSTRAINT accounts_pk PRIMARY KEY (id)
);

ALTER TABLE public.accounts ADD CONSTRAINT accounts_id_subject_fkey FOREIGN KEY (id_subject) REFERENCES subjects(id);



-- Объекты

CREATE TABLE public.objects (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	object_description text NOT NULL,
	object_address text NOT NULL,
	object_owner int8 NOT NULL,
	CONSTRAINT objects_pk PRIMARY KEY (id)
);


-- public.objects foreign keys

ALTER TABLE public.objects ADD CONSTRAINT objects_owner_fkey FOREIGN KEY (object_owner) REFERENCES subjects(id);



-- Изображения объектов

CREATE TABLE public.pics (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	id_object bigint NOT NULL,
	picture_url text NOT NULL,
	picture_description text NULL,
	CONSTRAINT pics_pk PRIMARY KEY (id)
);


-- public.pics foreign keys

ALTER TABLE public.pics ADD CONSTRAINT pics_id_object_fkey FOREIGN KEY (id_object) REFERENCES objects(id);



-- Аренды

CREATE TABLE public.rents (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	id_subject bigint NULL,
	id_object bigint NULL,
	price money NOT NULL,
	start_time timestamp(0) NOT NULL,
	expire_time timestamp(0) NOT NULL,
	CONSTRAINT rents_pk PRIMARY KEY (id),
	CONSTRAINT rents_price_check CHECK (((price)::numeric >= 0.0))
);


-- public.rents foreign keys

ALTER TABLE public.rents ADD CONSTRAINT rents_id_object_fkey FOREIGN KEY (id_object) REFERENCES objects(id);
ALTER TABLE public.rents ADD CONSTRAINT rents_id_subject_fkey FOREIGN KEY (id_subject) REFERENCES subjects(id);