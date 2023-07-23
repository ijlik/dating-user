package repository

import (
	"context"
	"database/sql"
	"errors"
)

func (r *repo) GetProfileBySwiperId(
	ctx context.Context,
	swiperId string,
) ([]*Profile, error) {
	randomProfile1, err := r.getRandomProfile(ctx, swiperId)
	if err != nil {
		return nil, err
	}
	if randomProfile1 == nil {
		return nil, errors.New("no more profile to show")
	}

	randomProfile2, err := r.getProfileWithoutId(ctx, swiperId, randomProfile1.ID)
	if err != nil {
		return nil, err
	}
	err = r.CreateSwipes(ctx, &Swipe{
		SwiperId: swiperId,
		SwipedId: randomProfile1.ID,
		IsLike:   sql.NullBool{},
	})
	if err != nil {
		return nil, err
	}
	var result []*Profile
	result = append(result, randomProfile1)
	if randomProfile2 != nil {
		result = append(result, randomProfile2)
	}
	return result, nil
}

func (r *repo) GetProfileBySwiperIdWithProfileId(
	ctx context.Context,
	swiperId,
	profileId string,
) ([]*Profile, error) {
	currentProfile, err := r.getProfileById(ctx, swiperId, profileId)
	if err != nil {
		return nil, err
	}
	if currentProfile == nil {
		return nil, errors.New("not found")
	}

	randomProfile, err := r.getProfileWithoutId(ctx, swiperId, currentProfile.ID)
	if err != nil {
		return nil, err
	}

	err = r.CreateSwipes(ctx, &Swipe{
		SwiperId: swiperId,
		SwipedId: currentProfile.ID,
		IsLike:   sql.NullBool{},
	})
	if err != nil {
		return nil, err
	}

	var result []*Profile
	result = append(result, currentProfile)
	if randomProfile != nil {
		result = append(result, randomProfile)
	}
	return result, nil
}

const getRandomProfileQuery = `SELECT id, user_id, name, birth_date, gender, photos, hobby,	interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> $1 AND id NOT IN (SELECT swiped_id FROM swipes WHERE swiper_id = $1 AND DATE(created_at) = CURRENT_DATE) ORDER BY RANDOM() LIMIT 1`

func (r *repo) getRandomProfile(
	ctx context.Context,
	swiperId string,
) (*Profile, error) {
	var data Profile
	err := r.conn.GetContext(
		ctx,
		&data,
		getRandomProfileQuery,
		swiperId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

const getProfileByIdQuery = `SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> $1 AND id = $2 AND id NOT IN (SELECT swiped_id FROM swipes WHERE swiper_id = $1 AND DATE(created_at) = CURRENT_DATE) LIMIT 1`

func (r *repo) getProfileById(
	ctx context.Context,
	swiperId,
	profileId string,
) (*Profile, error) {
	var data Profile
	err := r.conn.GetContext(
		ctx,
		&data,
		getProfileByIdQuery,
		swiperId,
		profileId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

const getProfileWithoutIdQuery = `SELECT id, user_id, name, birth_date, gender, photos, hobby, interest, location, is_premium, is_premium_valid_until, daily_swap_quota, created_at, updated_at FROM profiles WHERE name <> '' AND birth_date < CURRENT_TIMESTAMP AND gender <> '' AND photos <> '' AND hobby <> '' AND interest <> '' AND location <> '' AND id <> $1 AND id <> $2 AND id NOT IN (SELECT swiped_id FROM swipes WHERE swiper_id = $1 AND DATE(created_at) = CURRENT_DATE) LIMIT 1`

func (r *repo) getProfileWithoutId(
	ctx context.Context,
	swiperId,
	profileId string,
) (*Profile, error) {
	var data Profile
	err := r.conn.GetContext(
		ctx,
		&data,
		getProfileWithoutIdQuery,
		swiperId,
		profileId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

const createSwipesQuery = `INSERT INTO swipes (swiper_id, swiped_id, is_like, created_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`

const checkIfHasSwipeQuery = `SELECT count(*) FROM swipes WHERE swiper_id = $1 AND swiped_id = $2 AND is_like IN (true, false) AND DATE(created_at) = CURRENT_DATE`

const deleteSwipesShowOnlyQuery = `DELETE FROM swipes WHERE swiper_id = $1 AND swiped_id = $2 AND DATE(created_at) = CURRENT_DATE`

func (r *repo) CreateSwipes(
	ctx context.Context,
	req *Swipe,
) error {
	var count int
	err := r.conn.GetContext(
		ctx,
		&count,
		checkIfHasSwipeQuery,
		req.SwiperId,
		req.SwipedId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			count = 0
		}

		count = 0
	}
	if count > 0 {
		return errors.New("you already swipe this profile today")
	}

	if _, err := r.conn.ExecContext(
		ctx,
		deleteSwipesShowOnlyQuery,
		req.SwiperId,
		req.SwipedId,
	); err != nil {
		return err
	}

	if _, err := r.conn.ExecContext(
		ctx,
		createSwipesQuery,
		req.RowData()...,
	); err != nil {
		return err
	}

	return nil
}

const getSwipesCountAttribute = `SELECT count(*) FROM swipes WHERE swiper_id = $1 AND DATE(created_at) = CURRENT_DATE`

func (r *repo) GetSwipesCount(
	ctx context.Context,
	swiperId string,
) (int, error) {
	var count int
	err := r.conn.GetContext(
		ctx,
		&count,
		getSwipesCountAttribute,
		swiperId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}
